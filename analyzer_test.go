package main

import "testing"

func TestIdentifyComment(t *testing.T) {
	tests := []struct {
		path, trackID, smell string
	}{
		{"a/b/c/ruby/d/e.md", "ruby", "d/e"},
		{"a/b/c/go/d.md", "go", "d"},
	}

	for _, test := range tests {
		trackID, smell := identifyComment("a/b/c", test.path)
		if trackID != test.trackID {
			t.Errorf("trackID - got: %s, want: %s", trackID, test.trackID)
		}

		if smell != test.smell {
			t.Errorf("smell - got: %s, want: %s", smell, test.smell)
		}
	}
}
