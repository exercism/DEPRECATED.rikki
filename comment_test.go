package main

import "testing"

const fakeCommentDir = "fixtures/feedback"

func TestExistingComment(t *testing.T) {
	c := NewAnalyzerComment(fakeCommentDir, "ruby", "loops", "nesting")
	b, err := c.Bytes()
	if err != nil {
		t.Errorf("Unexpected error in comment at %s - %s", c.path, err)
	}
	expected := "Nope.\n"
	if string(b) != expected {
		t.Errorf("Expected %s - Got %s", expected, string(b))
	}
}

func TestMissingComment(t *testing.T) {
	c := NewAnalyzerComment(fakeCommentDir, "ruby", "loops", "too-many")
	b, err := c.Bytes()
	if err == nil {
		t.Errorf("Comment should be missing at %s", c.path)
	}
	if len(b) != 0 {
		t.Error("Comment should be an empty string when the problem is missing.")
	}
}
