package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/exercism/rikki/analysis/crystal"
	"github.com/exercism/rikki/analysis/golang"
	"github.com/exercism/rikki/analysis/ruby"
	"github.com/jrallison/go-workers"
)

// Analyzer is a job that provides feedback on specific issues in the code.
// The job receives the uuid of a submission, calls the exercism API to get
// the code, submits the code for static analysis, and then,
// based on the results, chooses a response to submit as a comment from rikki-
// back to the conversation on exercism.
type Analyzer struct {
	exercism *Exercism
	comments map[string]map[string][]byte
}

type analyzeFunc func(map[string]string) ([]string, error)

// NewAnalyzer configures an analyzer job to talk to the exercism and whatever analysis APIs we're using.
// We load the comments from disc when we create the analyzer.
// This means that rikki- has to be restarted if we update the comments.
func NewAnalyzer(exercism *Exercism, dir string) (*Analyzer, error) {
	dir = filepath.Join(dir, "analyzer")

	comments := make(map[string]map[string][]byte)

	fn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		b, err := read(path)
		if err != nil {
			return err
		}
		trackID, smell := identifyComment(dir, path)
		if comments[trackID] == nil {
			comments[trackID] = make(map[string][]byte)
		}
		comments[trackID][smell] = b

		return nil
	}

	if err := filepath.Walk(dir, fn); err != nil {
		return nil, err
	}

	return &Analyzer{
		exercism: exercism,
		comments: comments,
	}, nil
}

func identifyComment(dir, path string) (trackID, smell string) {
	r := strings.NewReplacer(dir, "", ".md", "")
	path = r.Replace(path)

	segments := strings.Split(path, string(filepath.Separator))

	if len(segments) < 2 {
		return "", ""
	}

	return segments[1], strings.Join(segments[2:], string(filepath.Separator))
}

func (analyzer *Analyzer) process(msg *workers.Msg) {
	// Fetch the solution from the Exercism API.
	uuid, err := msg.Args().GetIndex(0).String()
	if err != nil {
		lgr.Printf("unable to determine submission key - %s\n", err)
		return
	}
	solution, err := analyzer.exercism.FetchSolution(uuid)
	if err != nil {
		lgr.Printf("%s\n", err)
		return
	}

	// Detect known smells.
	var fn analyzeFunc
	switch solution.TrackID {
	case "ruby":
		fn = ruby.Analyze
	case "go":
		fn = golang.Analyze
	case "crystal":
		fn = crystal.Analyze
	default:
		lgr.Printf("skipping - rikki- doesn't support %s\n", solution.TrackID)
		return
	}
	smells, err := fn(solution.Files)
	if err != nil {
		lgr.Printf("%s - %s", uuid, err)
		return
	}

	// Log what we found.
	sanity := log.New(os.Stdout, "SANITY: ", log.Ldate|log.Ltime|log.Lshortfile)
	for _, smell := range smells {
		sanity.Printf("%s : %s\n", uuid, smell)
	}

	// Select the first smell that we have a comment for.
	var comment []byte
	for _, smell := range smells {
		b := analyzer.comments[solution.TrackID][smell]

		if len(b) > 0 {
			comment = b
			break
		}
	}
	if len(comment) == 0 {
		return
	}

	// Submit the comment back to the Exercism API.
	if err := analyzer.exercism.SubmitComment(comment, uuid); err != nil {
		lgr.Printf("%s\n", err)
	}
}
