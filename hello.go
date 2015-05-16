package main

import "github.com/jrallison/go-workers"

// Hello is a job that provides encouragement after someone submits "Hello World".
// The job receives the uuid of a submission and submits a comment from rikki-
// to the conversation on exercism.
type Hello struct {
	exercism *Exercism
}

// NewHello configures a Hello job to talk to the exercism API.
func NewHello(exercism *Exercism) *Hello {
	return &Hello{
		exercism: exercism,
	}
}

func (hello *Hello) process(msg *workers.Msg) {
	submissionUuid, err := msg.Args().GetIndex(0).String()
	if err != nil {
		lgr.Printf("unable to determine submission uuid - %s\n", err)
		return
	}

	// load rikki-'s encouragement
	c := NewHelloComment("")
	body, err := c.Bytes()
	if err != nil {
		lgr.Printf("%s\n", err)
		return
	}

	if err := hello.exercism.SubmitComment(body, submissionUuid); err != nil {
		lgr.Printf("%s\n", err)
	}
}
