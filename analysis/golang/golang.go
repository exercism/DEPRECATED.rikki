package golang

import (
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	smellFmt = "gofmt"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type solution struct {
	files map[string]string
	dir   string
}

func (s *solution) write() error {
	s.dir = path.Join(os.TempDir(), strconv.Itoa(rand.Intn(1e9)))
	os.Mkdir(s.dir, os.ModePerm)

	for name, code := range s.files {
		filename := path.Join(s.dir, name)

		// Fix any potential trailing newline issues.
		// These wouldn't be visible to a reviewer, and if they're not
		// running gofmt, then eventually we'll catch it with a more obvious problem.

		code = strings.TrimRight(code, "\n") + "\n"

		if err := ioutil.WriteFile(filename, []byte(code), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// Analyze detects certain issues in Go code.
func Analyze(files map[string]string) ([]string, error) {
	s := &solution{
		files: files,
	}
	if err := s.write(); err != nil {
		return nil, err
	}
	defer os.Remove(s.dir)

	ok, err := isGofmted(s)
	if err != nil {
		return nil, err
	}
	if !ok {
		return []string{smellFmt}, nil
	}
	return nil, nil
}

func isGofmted(s *solution) (bool, error) {
	output, err := exec.Command("gofmt", "-l", s.dir).Output()
	if err != nil {
		return false, err
	}
	lines := strings.Split(string(output), "\n")

	for name := range s.files {
		for _, line := range lines {
			if strings.HasSuffix(line, name) {
				return false, nil
			}
		}
	}
	return true, nil
}
