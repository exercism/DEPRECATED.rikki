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
	smellVet   = `go-vet`
	smellStub  = `stub`
	smellBuild = `build-constraint`
	smellCase  = `mixed-caps`
	smellZero  = `zero-value`

	msgAllCaps   = `don't use ALL_CAPS in Go names`
	msgSnakeCase = `don't use underscores in Go names`
)

var (
	rgxStub = regexp.MustCompile(`\bstub\b`)
	rgxZero = regexp.MustCompile(`should drop.*from declaration of.*it is the zero value`)
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
		smellVet:   isVetted,
		smellStub:  isStubless,
		smellBuild: noBuildConstraint,
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
	linted, err := lint(s)
	if err != nil {
		if len(smells) > 0 {
			return smells, nil
		}
		return nil, err
	}
	smells = append(smells, linted...)

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

func isVetted(s *solution) (bool, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return false, err
	}
	defer os.Chdir(pwd)

	os.Chdir(s.dir)

	output, _ := exec.Command("go", "vet", `./...`).CombinedOutput()
	return string(output) == "", nil
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

func lint(s *solution) ([]string, error) {
	output, err := exec.Command("golint", filepath.Join(s.dir, `...`)).CombinedOutput()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	m := map[string]bool{}
	for _, line := range lines {
		if isMixedCaps(line) {
			m[smellCase] = true
		}
		if isZeroValue(line) {
			m[smellZero] = true
		}
	}
	var smells []string
	for smell := range m {
		smells = append(smells, smell)
	}

	return smells, nil
}

func isMixedCaps(msg string) bool {
	return strings.Contains(msg, msgSnakeCase) || strings.Contains(msg, msgAllCaps)
}

func isZeroValue(msg string) bool {
	return rgxZero.Match([]byte(msg))
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
