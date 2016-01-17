package golang

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type solution struct {
	files map[string]string
	dir   string
}

func newSolution(m map[string]string) *solution {
	files := map[string]string{}
	for name, code := range m {
		// Normalize path names.
		name = strings.Replace(name, `\`, string(filepath.Separator), -1)
		name = strings.Replace(name, `/`, string(filepath.Separator), -1)
		name = `/` + strings.TrimLeft(name, `/`)

		// Fix any potential trailing newline issues.
		// These wouldn't be visible to a reviewer, and if they're not
		// running gofmt, then eventually we'll catch it with a more obvious problem.
		files[name] = strings.TrimRight(code, "\n") + "\n"
	}

	return &solution{
		files: files,
	}
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