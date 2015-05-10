package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Comment struct {
	path string
}

func NewAnalyzerComment(dir, language, category, issue string) *Comment {
	if dir == "" {
		dir = commentDir()
	}
	return &Comment{
		path: fmt.Sprintf("%s/analyzer/%s/%s/%s.md", dir, language, category, issue),
	}
}

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
