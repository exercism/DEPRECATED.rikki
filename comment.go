package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Comment is an observation by rikki-.
// It will be submitted to a specific submission on exercism.
type Comment struct {
	path string
}

// NewAnalyzerComment creates a comment that provides critical feedback on some code.
func NewAnalyzerComment(dir, language, category, issue string) *Comment {
	if dir == "" {
		dir = commentDir()
	}
	return &Comment{
		path: fmt.Sprintf("%s/analyzer/%s/%s/%s.md", dir, language, category, issue),
	}
}

// NewHelloComment creates a comment that provides generic encouragement.
// The "hello world" exercise is more about ensuring that everything is wired up
// than about good, idiomatic code, and this comment reflects that purpose.
func NewHelloComment(dir string) *Comment {
	if dir == "" {
		dir = commentDir()
	}
	return &Comment{
		path: fmt.Sprintf("%s/hello/hello.md", dir),
	}
}

// Bytes returns the actual text of the comment.
func (c *Comment) Bytes() ([]byte, error) {
	var b []byte
	if _, err := os.Stat(c.path); err != nil {
		return b, err
	}
	comment, err := ioutil.ReadFile(c.path)
	if err != nil {
		return b, err
	}
	return comment, nil
}
