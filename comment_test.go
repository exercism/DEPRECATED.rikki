package main

import "testing"

const fakeCommentDir = "fixtures/feedback"

func TestExistingComment(t *testing.T) {
	c := NewComment("ruby", "loops", "nesting")
	c.Dir = fakeCommentDir
	b, err := c.Bytes()
	if err != nil {
		t.Errorf("Unexpected error in comment at %s - %s", c.Path(), err)
	}
	expected := "Nope.\n"
	if string(b) != expected {
		t.Errorf("Expected %s - Got %s", expected, string(b))
	}
}

func TestMissingComment(t *testing.T) {
	c := NewComment("ruby", "loops", "too-many")
	c.Dir = fakeCommentDir

	b, err := c.Bytes()
	if err == nil {
		t.Errorf("Comment should be missing at %s", c.Path())
	}
	if len(b) != 0 {
		t.Error("Comment should be an empty string when the problem is missing.")
	}
}
