package golang

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type solution struct {
	files    map[string]string
	comments []string
	dir      string
}

func newSolution(m map[string]string) *solution {
	files := map[string]string{}
	for name, code := range m {
		// Normalize path names.
		name = strings.Replace(name, `\`, string(filepath.Separator), -1)
		name = strings.Replace(name, `/`, string(filepath.Separator), -1)
		name = `/` + strings.TrimLeft(name, `/`)

		files[name] = normalizeSource(code)
	}

	return &solution{
		files: files,
	}
}

// Fix potential issues with leading or trailing newlines.
// These wouldn't be visible to a reviewer, and if they're not
// running gofmt, then eventually we'll catch it with a more obvious problem.
func normalizeSource(code string) string {
	code = strings.TrimRight(code, "\n") + "\n"
	code = strings.TrimLeft(code, "\n")
	return strings.Replace(code, "\r\n", "\n", -1)
}

func (s *solution) write() error {
	s.dir = path.Join(os.TempDir(), strconv.Itoa(rand.Intn(1e9)))
	if err := os.Mkdir(s.dir, os.ModePerm); err != nil {
		return err
	}

	for name, code := range s.files {
		filename := path.Join(s.dir, name)

		if err := os.MkdirAll(path.Dir(filename), os.ModePerm); err != nil {
			return err
		}

		if err := ioutil.WriteFile(filename, []byte(code), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (s *solution) extractComments() error {
	for name, code := range s.files {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, name, code, parser.ParseComments)
		if err != nil {
			return err
		}
		for _, cg := range f.Comments {
			s.comments = append(s.comments, cg.Text())
		}
	}
	return nil
}
