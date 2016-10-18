package crystal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Host is the base URL for the analysseur API.
var Host string

type request struct {
	ID       string `json:"id"`
	Contents string `json:"contents"`
}

type response struct {
	ID       string    `json:"id"`
	Problems []problem `json:"problems"`
	Error    string    `json:"error"`
}

type problem struct {
	Type   string `json:"type"`
	Result string `json:"result"`
}

// Analyze Crystal code for formatting errors (and other bad things maybe later)
func Analyze(files map[string]string) ([]string, error) {
	var sources []string
	for _, source := range files {
		sources = append(sources, source)
	}

	url := fmt.Sprintf("%s/check", Host)
	code := strings.Join(sources, "\n")
	requestBody := request{ID: "rikki", Contents: code}
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBodyJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s responded with status %d - %s\n", url, resp.StatusCode, string(respBody))
	}

	var res response
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return nil, err
	}

	if res.Error != "" {
		return nil, errors.New(res.Error)
	}

	var smells []string
	for _, prob := range res.Problems {
		if prob.Result == "true" {
			smells = append(smells, prob.Type)
		}
	}

	return smells, nil
}
