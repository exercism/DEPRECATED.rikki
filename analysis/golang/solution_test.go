package golang

import (
	"os"
	"sort"
	"testing"
)

func TestNewSolution(t *testing.T) {
	files := map[string]string{
		`some/subdir/code.go`: "",
		`/other/dir/code.go`:  "",
		`\win\dows\code.go`:   "",
	}
	s := newSolution(files)

	want := []string{
		`/other/dir/code.go`,
		`/some/subdir/code.go`,
		`/win/dows/code.go`,
	}
	got := []string{}
	for name := range s.files {
		got = append(got, name)
	}
	sort.Strings(got)

	for i := 0; i < 3; i++ {
		if want[i] != got[i] {
			t.Errorf("got %s, want %s", got[i], want[i])
		}
	}
}

func TestWrite(t *testing.T) {
	files := map[string]string{
		`code.go`:            "",
		`some/code.go`:       "",
		`/other/dir/code.go`: "",
	}
	s := newSolution(files)
	if err := s.write(); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(s.dir)

	for filename := range s.files {
		if _, err := os.Stat(s.dir + filename); err != nil {
			if os.IsNotExist(err) {
				t.Error(err)
			} else {
				t.Fatal(err)
			}
		}
	}
}
