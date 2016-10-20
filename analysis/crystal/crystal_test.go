package crystal

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
)

var (
	u string
)

func init() {
	Host = "http://localhost:3000"
	u = fmt.Sprintf("%s/check", Host)
}

func TestFormattedCode(t *testing.T) {
	mockJSON := `{"id":"rikki", "problems":[{"type":"unformatted", "result":"false"}], "error":""}`
	result, err := getTestResponse(mockJSON)
	if err != nil {
		t.Fatal(err)
	}

	actual := len(result)
	expected := 0
	if actual != expected {
		t.Errorf("got %v, want %v", actual, expected)
	}
}

func TestUnformattedCode(t *testing.T) {
	mockJSON := `{"id":"rikki", "problems":[{"type":"unformatted", "result":"true"}], "error":""}`
	result, err := getTestResponse(mockJSON)
	if err != nil {
		t.Fatal(err)
	}

	actual := result[0]
	expected := "unformatted"
	if actual != expected {
		t.Errorf("got %v, want %v", actual, expected)
	}
}

func getTestResponse(mockJSON string) ([]string, error) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockRes := httpmock.NewStringResponder(200, mockJSON)
	httpmock.RegisterResponder("POST", u, mockRes)

	var files map[string]string
	files = make(map[string]string)
	files["test.cr"] = "code"

	result, err := Analyze(files)
	if err != nil {
		return nil, err
	}

	return result, nil
}
