package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Exercism is a client that talks to the exercism API.
type Exercism struct {
	Host string
	Auth string
}

// NewExercism creates an exercism client, configured to talk to the API.
func NewExercism(host, auth string) *Exercism {
	return &Exercism{Host: host, Auth: auth}
}

type codePayload struct {
	TrackID       string            `json:"track_id"`
	SolutionFiles map[string]string `json:"solution_files"`
	Error         string            `json:"error"`
}

type commentBody struct {
	Comment string `json:"comment"`
}

// Solution is an iteration of a specific problem in a particular language.
type Solution struct {
	TrackID string
	Files   map[string]string
}

// FetchSolution fetches the code of a solution from the exercism API.
func (e *Exercism) FetchSolution(uuid string) (*Solution, error) {
	url := fmt.Sprintf("%s/api/v1/submissions/%s", e.Host, uuid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare request to %s - %s\n", url, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed - %s\n", url, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response - %s\n", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s responded with status %d - %s\n", url, resp.StatusCode, string(body))
	}

	var cp codePayload
	if err := json.Unmarshal(body, &cp); err != nil {
		return nil, fmt.Errorf("%s - %s\n", uuid, err)
	}

	return &Solution{TrackID: cp.TrackID, Files: cp.SolutionFiles}, nil
}

// SubmitComment submits a rikki- comment to a particular submission via the exercism API.
func (e *Exercism) SubmitComment(comment []byte, uuid string) error {
	experiment := "_This is an automated review based on lots and lots of real-life reviews. [Read more](http://exercism.io/rikki) about this experiment._"
	s := fmt.Sprintf("%s\n-----\n%s", string(comment), experiment)

	cb, err := json.Marshal(&commentBody{Comment: s})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/submissions/%s/comments?shared_key=%s", e.Host, uuid, e.Auth)
	req, err := http.NewRequest("POST", url, bytes.NewReader(cb))
	if err != nil {
		return fmt.Errorf("cannot prepare request to %s - %s", url, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request to %s failed - %s", url, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("%s responded with status %d - %s", url, resp.StatusCode, string(body))
	}
	return nil
}
