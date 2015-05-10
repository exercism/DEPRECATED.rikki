package main

import "github.com/jrallison/go-workers"

type Hello struct {
	exercism *Exercism
}

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
