package golang

import (
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/golang/lint"
)

const (
	smellFmt   = `gofmt`
	smellVet   = `go-vet`
	smellStub  = `stub`
	smellBuild = `build-constraint`
	smellCase  = `mixed-caps`
	smellZero  = `zero-value`
	smellElse  = `if-return-else`

	msgAllCaps   = `don't use ALL_CAPS in Go names`
	msgSnakeCase = `don't use underscores in Go names`
	msgOutdent   = `if block ends with a return statement, so drop this else and outdent its block`
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

	if err := s.extractComments(); err != nil {
		return nil, err
	}

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
	linted, err := lintify(s)
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
	for _, c := range s.comments {
		if rgxStub.Match([]byte(c)) {
			return false, nil
		}
	}
	return true, nil
}

func noBuildConstraint(s *solution) (bool, error) {
	for _, c := range s.comments {
		if strings.Contains(c, `+build !example`) {
			return false, nil
		}
	}
	return true, nil
}

func lintify(s *solution) ([]string, error) {
	linter := &lint.Linter{}
	m := map[string]bool{}

	for filename, src := range s.files {
		problems, err := linter.Lint(filename, []byte(src))
		if err != nil {
			return nil, err
		}
		for _, problem := range problems {
			if problem.Category == "zero-value" {
				m[smellZero] = true
			}
			if problem.Category == "indent" {
				m[smellElse] = true
			}
			if problem.Category == "naming" && isMixedCaps(problem.Text) {
				m[smellCase] = true
			}
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
