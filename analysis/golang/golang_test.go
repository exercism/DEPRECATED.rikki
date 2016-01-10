package golang

import (
	"os"
	"testing"
)

var codeBad = `package bad

func ok() {
	println(3 % 2 == 0)
}
`

var codeGood = `package good

func ok() {
	println(3%2 == 0)
}
`

var codeNewline = `package nl

func ok() {
	println(3%2 == 0)
}`

var tests = []struct {
	desc, code string
	ok         bool
}{
	{"good", codeGood, true},
	{"bad", codeBad, false},
	{"newline", codeNewline, true},
}

func TestGofmted(t *testing.T) {
	for _, test := range tests {
		s := &solution{
			files: map[string]string{"code.go": test.code},
		}
		if err := s.write(); err != nil {
			t.Fatal(err)
		}
		defer os.Remove(s.dir)

		ok, err := isGofmted(s)
		if err != nil {
			t.Fatal(err)
		}
		if ok != test.ok {
			t.Errorf("%s: got %t, want %t", test.desc, ok, !ok)
		}
	}
}

func TestAnalyze(t *testing.T) {
	for _, test := range tests {
		smells, err := Analyze(map[string]string{"code.go": test.code})
		if err != nil {
			t.Fatal(err)
		}
		if test.ok && len(smells) != 0 {
			t.Errorf("%s: got %v, want empty list", test.desc, smells)
		}
		if !test.ok && len(smells) == 0 {
			t.Errorf("%s: got empty list, want 'gofmt'", test.desc)
		}
	}
}
