package crystal

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnalyze(t *testing.T) {
	tests := []struct {
		payload string
		smells  []string
	}{
		{
			`{"id":"rikki", "problems":[{"type":"unformatted", "result":"false"}], "error":""}`,
			[]string{},
		},
		{
			`{"id":"rikki", "problems":[{"type":"unformatted", "result":"true"}], "error":""}`,
			[]string{"unformatted"},
		},
	}

	for _, test := range tests {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(test.payload))
		}))
		defer ts.Close()

		// Fake out the host and API endpoint.
		Host = ts.URL
		Path = ""

		// Fake out the files to analyze. We only care what the server responds.
		smells, err := Analyze("", map[string]string{"test.cr": "code"})
		if err != nil {
			t.Fatal(err)
		}

		if len(smells) != len(test.smells) {
			t.Errorf("Got %d smells, expected %d", len(smells), len(test.smells))
			continue
		}

		for i, smell := range smells {
			if test.smells[i] != smell {
				t.Errorf("Got smell %s at index %d, expected %s.", smell, i, test.smells[i])
			}
		}
	}
}
