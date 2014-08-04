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

type analysisBody struct {
	Code string `json:"code"`
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
		lgr.Printf("%s responded with status %v - %v", url, resp.StatusCode, string(body))
		return
	}

	var cp codePayload
	err = json.Unmarshal(body, &cp)
	if err != nil {
		lgr.Println(err)
		return
	}

	if cp.Language != "ruby" {
		lgr.Println("Skipping code in %s", cp.Language)
		return
	}

	// Step 2: submit code to analysseur
	url = fmt.Sprintf("%s/analyze/%s", analyzer.analysseurHost, cp.Language)
	ab := analysisBody{Code: cp.Code}
	abJSON, err := json.Marshal(ab)
	if err != nil {
		lgr.Println(err)
		return
	}

	req, err = http.NewRequest("POST", url, bytes.NewReader(abJSON))
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
	if resp.StatusCode != http.StatusOK {
		lgr.Printf("%s responded with status %v - %v", url, resp.StatusCode, string(body))
		return
	}

	var ap analysisPayload
	err = json.Unmarshal(body, &ap)
	if err != nil {
		lgr.Println(err)
		return
	}

	if ap.Error != "" {
		lgr.Printf("analysis api is complaining - %s", ap.Error)
		return
	}

	if len(ap.Results) == 0 {
		// no feedback, bailing
		return
	}

	// Step 3: Find comments based on analysis result
	var comments [][]byte
	for _, result := range ap.Results {
		for _, key := range result.Keys {
			filename := fmt.Sprintf("comments/%s/%s/%s.md", cp.Language, result.Type, key)
			if _, err = os.Stat(filename); err != nil {
				if os.IsNotExist(err) {
					lgr.Printf("we need to add a comment for %s - %s\n", result.Type, key)
					return
				}
				lgr.Println(err)
			}
			comment, err := ioutil.ReadFile(filename)
			if err != nil {
				lgr.Printf("unable to read %s - %v", filename, err)
			} else {
				comments = append(comments, comment)
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
