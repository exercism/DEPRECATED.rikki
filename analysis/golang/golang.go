package golang

import (
	"go/parser"
	"go/token"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	smellFmt   = `gofmt`
	smellStub  = `stub`
	smellBuild = `build-constraint`
	smellCase  = `mixed-caps`

	msgAllCaps   = `don't use ALL_CAPS in Go names`
	msgSnakeCase = `don't use underscores in Go names`
)

var (
	rgxStub = regexp.MustCompile(`\bstub\b`)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Analyze detects certain issues in Go code.
func Analyze(files map[string]string) ([]string, error) {
	s := newSolution(files)
	if err := s.write(); err != nil {
		return nil, err
	}
	defer os.Remove(s.dir)

	smells := []string{}

	detectors := map[string]func(*solution) (bool, error){
		smellFmt:   isGofmted,
		smellStub:  isStubless,
		smellBuild: noBuildConstraint,
		smellCase:  usesMixedCaps,
	}

	for smell, fn := range detectors {
		ok, err := fn(s)
		if err != nil {
			return nil, err
		}
		if !ok {
			smells = append(smells, smell)
		}
	}

	return smells, nil
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

func isStubless(s *solution) (bool, error) {
	comments, err := astComments(s)
	if err != nil {
		return false, err
	}
	for _, c := range comments {
		if rgxStub.Match([]byte(c)) {
			return false, nil
		}
	}
	return true, nil
}

func noBuildConstraint(s *solution) (bool, error) {
	comments, err := astComments(s)
	if err != nil {
		return false, err
	}
	for _, c := range comments {
		if strings.Contains(c, `+build !example`) {
			return false, nil
		}
	}
	return true, nil
}

func usesMixedCaps(s *solution) (bool, error) {
	output, err := exec.Command("golint", filepath.Join(s.dir, `...`)).CombinedOutput()
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, msgSnakeCase) || strings.Contains(line, msgAllCaps) {
			return false, nil
		}
	}
	return true, nil
}

func astComments(s *solution) ([]string, error) {
	comments := []string{}
	for name, code := range s.files {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, name, code, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		for _, cg := range f.Comments {
			comments = append(comments, cg.Text())
		}
	}
	return comments, nil
}
