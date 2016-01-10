package ruby

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
)

// Host is the base URL for the analysseur API.
var Host string

type result struct {
	Type string   `json:"type"`
	Keys []string `json:"keys"`
}
type payload struct {
	Results []result `json:"results"`
	Error   string   `json:"error"`
}

// Analyze detects a specific set of code smells in Ruby code.
func Analyze(files map[string]string) ([]string, error) {
	var sources []string
	for _, source := range files {
		sources = append(sources, source)
	}

	// Step 2: submit code to analysseur
	url := fmt.Sprintf("%s/analyze/ruby", Host)
	codeBody := struct {
		Code string `json:"code"`
	}{
		strings.Join(sources, "\n"),
	}
	codeBodyJSON, err := json.Marshal(codeBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(codeBodyJSON))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s responded with status %d - %s\n", url, resp.StatusCode, string(body))
	}

	var pld payload
	err = json.Unmarshal(body, &pld)
	if err != nil {
		return nil, err
	}

	if pld.Error != "" {
		return nil, errors.New(pld.Error)
	}

	var smells []string
	for _, result := range pld.Results {
		for _, key := range result.Keys {
			smells = append(smells, filepath.Join(result.Type, key))
		}
	}

	// shuffle code smells
	for i := range smells {
		j := rand.Intn(i + 1)
		smells[i], smells[j] = smells[j], smells[i]
	}
	return smells, nil
}
