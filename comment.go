package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Comment struct {
	Dir      string
	Language string
	Category string
	Issue    string
}

func NewComment(language, category, issue string) *Comment {
	return &Comment{
		Dir:      commentDir(),
		Language: language,
		Category: category,
		Issue:    issue,
	}
}

func (c *Comment) Path() string {
	return fmt.Sprintf("%s/%s/%s/%s.md", c.Dir, c.Language, c.Category, c.Issue)
}

func (c *Comment) Bytes() ([]byte, error) {
	var b []byte
	if _, err := os.Stat(c.Path()); err != nil {
		return b, err
	}
	comment, err := ioutil.ReadFile(c.Path())
	if err != nil {
		return b, err
	}
	return comment, nil
}

func commentDir() string {
	dir := os.Getenv("RIKKI_FEEDBACK_DIR")
	if dir == "" {
		dir = "comments"
	}
	return dir
}
