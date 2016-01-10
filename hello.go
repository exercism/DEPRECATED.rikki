package main

import "github.com/jrallison/go-workers"

// Hello is a job that provides encouragement after someone submits "Hello World".
// The job receives the uuid of a submission and submits a comment from rikki-
// to the conversation on exercism.
type Hello struct {
	exercism *Exercism
	comment  []byte
}

// NewHello configures a Hello job to talk to the exercism API.
func NewHello(exercism *Exercism) (*Hello, error) {
	b, err := NewHelloComment("").Bytes()
	if err != nil {
		return nil, err
	}
	return &Hello{
		exercism: exercism,
		comment:  b,
	}, nil
}

func (hello *Hello) process(msg *workers.Msg) {
	uuid, err := msg.Args().GetIndex(0).String()
	if err != nil {
		lgr.Printf("unable to determine submission uuid - %s\n", err)
		return
	}

	if err := hello.exercism.SubmitComment(hello.comment, uuid); err != nil {
		lgr.Printf("%s\n", err)
	}
}
