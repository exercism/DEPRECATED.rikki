package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/jrallison/go-workers"
)

type Analyzer struct {
	exercism       *Exercism
	analysseurHost string
}

func NewAnalyzer(exercism *Exercism, analysseur string) *Analyzer {
	return &Analyzer{
		exercism:       exercism,
		analysseurHost: analysseur,
	}
}

type analysisResult struct {
	Type string   `json:"type"`
	Keys []string `json:"keys"`
}
type analysisPayload struct {
	Results []analysisResult `json:"results"`
	Error   string           `json:"error"`
}

func (analyzer *Analyzer) process(msg *workers.Msg) {
	submissionUuid, err := msg.Args().GetIndex(0).String()
	if err != nil {
		lgr.Printf("unable to determine submission key - %s\n", err)
		return
	}

	solution, err := analyzer.exercism.FetchSolution(submissionUuid)
	if err != nil {
		lgr.Printf("%s\n", err)
		return
	}

	if solution.TrackID != "ruby" {
		lgr.Printf("skipping - rikki- doesn't support %s\n", solution.TrackID)
		return
	}

	// Step 2: submit code to analysseur
	url := fmt.Sprintf("%s/analyze/%s", analyzer.analysseurHost, solution.TrackID)
	codeBody := struct {
		Code string `json:"code"`
	}{
		strings.Join(solution.Sources, "\n"),
	}
	codeBodyJSON, err := json.Marshal(codeBody)
	if err != nil {
		lgr.Printf("%s - %s\n", submissionUuid, err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(codeBodyJSON))
	if err != nil {
		lgr.Printf("%s - cannot prepare request to %s - %s\n", submissionUuid, url, err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		lgr.Printf("%s - request to %s failed - %s\n", submissionUuid, url, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		lgr.Printf("%s - failed to fetch submission - %s\n", submissionUuid, err)
		return
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		lgr.Printf("%s - %s responded with status %d - %s\n", submissionUuid, url, resp.StatusCode, string(body))
		return
	}

	var ap analysisPayload
	err = json.Unmarshal(body, &ap)
	if err != nil {
		lgr.Printf("%s - %s\n", submissionUuid, err)
		return
	}

	if ap.Error != "" {
		lgr.Printf("analysis api is complaining about %s - %s\n", submissionUuid, ap.Error)
		return
	}

	if len(ap.Results) == 0 {
		// no feedback, bailing
		return
	}

	sanity := log.New(os.Stdout, "SANITY: ", log.Ldate|log.Ltime|log.Lshortfile)
	for _, result := range ap.Results {
		for _, key := range result.Keys {
			sanity.Printf("%s : %s - %s\n", submissionUuid, result.Type, key)
		}
	}

	// Step 3: Find comments based on analysis result
	// We are loading the results before choosing a comment at random
	// since not all results will have an associated comment, and it's
	// better to be a bit wasteful than to not submit a comment when
	// we could have.
	var comments [][]byte
	for _, result := range ap.Results {
		for _, key := range result.Keys {
			c := NewAnalyzerComment("", solution.TrackID, result.Type, key)
			b, err := c.Bytes()
			if err != nil {
				lgr.Printf("We probably need to add a comment at %s - %s\n", c.path, err)
			}
			if len(b) > 0 {
				comments = append(comments, b)
			}
		}
	}
	if len(comments) == 0 {
		// no comments, bailing
		return
	}

	// Step 4: submit random comment to exercism.io api
	comment := comments[rand.Intn(len(comments))]
	if err := analyzer.exercism.SubmitComment(comment, submissionUuid); err != nil {
		lgr.Printf("%s\n", err)
	}
}
