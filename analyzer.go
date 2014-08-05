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

	"github.com/jrallison/go-workers"
)

type Analyzer struct {
	exercismHost   string
	analysseurHost string
	auth           string
	Logger         *log.Logger
}

func NewAnalyzer(exercism, analysseur, auth string) *Analyzer {
	return &Analyzer{
		exercismHost:   exercism,
		analysseurHost: analysseur,
		auth:           auth,
		Logger:         log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

type codePayload struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Error    string `json:"error"`
}

type analysisResult struct {
	Type string   `json:"type"`
	Keys []string `json:"keys"`
}
type analysisPayload struct {
	Results []analysisResult `json:"results"`
	Error   string           `json:"error"`
}

type commentBody struct {
	Comment string `json:"comment"`
}

func (analyzer *Analyzer) process(msg *workers.Msg) {
	lgr := analyzer.Logger

	submissionKey, err := msg.Args().GetIndex(0).String()
	if err != nil {
		lgr.Printf("unable to determine submission key - %v\n", err)
		return
	}

	// Step 1: get code + language from exercism.io api
	url := fmt.Sprintf("%s/api/v1/submissions/%s", analyzer.exercismHost, submissionKey)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		lgr.Printf("cannot prepare request to %s - %v\n", url, err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		lgr.Printf("request to %s failed - %v\n", url, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		lgr.Printf("%s responded with status %v - %v\n", url, resp.StatusCode, string(body))
		return
	}

	var cp codePayload
	err = json.Unmarshal(body, &cp)
	if err != nil {
		lgr.Printf("%s - %v\n", submissionKey, err)
		return
	}

	if cp.Language != "ruby" {
		lgr.Println("Skipping code in %s", cp.Language)
		return
	}

	// Step 2: submit code to analysseur
	url = fmt.Sprintf("%s/analyze/%s", analyzer.analysseurHost, cp.Language)
	reqBody := struct {
		Code string `json:"code"`
	}{
		cp.Code,
	}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		lgr.Printf("%s - %v", submissionKey, err)
		return
	}

	req, err = http.NewRequest("POST", url, bytes.NewReader(reqBodyJSON))
	if err != nil {
		lgr.Printf("%s - cannot prepare request to %s - %v\n", submissionKey, url, err)
		return
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		lgr.Printf("%s - request to %s failed - %v\n", submissionKey, url, err)
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		lgr.Printf("%s - %s responded with status %v - %v", submissionKey, url, resp.StatusCode, string(body))
		return
	}

	var ap analysisPayload
	err = json.Unmarshal(body, &ap)
	if err != nil {
		lgr.Printf("%s - %v", submissionKey, err)
		return
	}

	if ap.Error != "" {
		lgr.Printf("analysis api is complaining about %s - %s", submissionKey, ap.Error)
		return
	}

	if len(ap.Results) == 0 {
		// no feedback, bailing
		return
	}

	sanity := log.New(os.Stdout, "SANITY: ", log.Ldate|log.Ltime|log.Lshortfile)
	for _, result := range ap.Results {
		for _, key := range result.Keys {
			sanity.Printf("%s : %s - %s\n", submissionKey, result.Type, key)
		}
	}

	// Step 3: Find comments based on analysis result
	var comments [][]byte
	for _, result := range ap.Results {
		for _, key := range result.Keys {
			c := NewComment(cp.Language, result.Type, key)
			b, err := c.Bytes()
			if err != nil {
				lgr.Printf("We probably need to add a comment at %s - %s\n", c.Path(), err)
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
	experiment := "_This is an automated nitpick. [Read more](http://exercism.io/rikki) about this experiment._"
	s := fmt.Sprintf("%s\n-----\n%s", string(comment), experiment)
	cb := commentBody{Comment: s}
	cbJSON, err := json.Marshal(cb)
	if err != nil {
		lgr.Println(err)
		return
	}

	url = fmt.Sprintf("%s/api/v1/submissions/%s/comments?shared_key=%s", analyzer.exercismHost, submissionKey, analyzer.auth)
	req, err = http.NewRequest("POST", url, bytes.NewReader(cbJSON))
	if err != nil {
		lgr.Printf("cannot prepare request to %s - %v\n", url, err)
		return
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		lgr.Printf("request to %s failed - %v\n", url, err)
		return
	}

	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		lgr.Printf("%s responded with status %v - %v", url, resp.StatusCode, string(body))
		return
	}
}
