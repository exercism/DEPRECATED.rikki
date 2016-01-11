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

var codeStub = `// This is a stub file
package main

// main is where everything begins.
func main() {
	println("ok")
}
`

var codeStubPkg = `package stub

func stub() bool {
	return true
}
`

var codeStubFragment = `// Package fragment has a ristubble in it.
// No matter that ristubble isn't actually a word.
package fragment

func ok() bool {
	return true
}
`

func TestGofmted(t *testing.T) {
	var tests = []struct {
		desc, code string
		ok         bool
	}{
		{"good", codeGood, true},
		{"bad", codeBad, false},
		{"newline", codeNewline, true},
	}

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

func TestStubs(t *testing.T) {
	var tests = []struct {
		desc, code string
		ok         bool
	}{
		{"comment", codeStub, false},
		{"pkg", codeStubPkg, true},
		{"fragment", codeStubFragment, true},
	}

	for _, test := range tests {
		s := &solution{
			files: map[string]string{test.desc + `.go`: test.code},
		}
		ok, err := isStubless(s)
		if err != nil {
			t.Fatal(err)
		}
		if ok != test.ok {
			t.Errorf("%s: got %t, want %t", test.desc, ok, !ok)
		}
	}
}

func TestAnalyze(t *testing.T) {
	var tests = []struct {
		desc, code string
		smells     []string
	}{
		{"good", codeGood, nil},
		{"bad", codeBad, []string{"gofmt"}},
		{"comment", codeStub, []string{"stub"}},
	}

	for _, test := range tests {
		smells, err := Analyze(map[string]string{"code.go": test.code})
		if err != nil {
			t.Fatal(err)
		}
		if len(test.smells) != len(smells) {
			t.Errorf("%s: got %v, want %v", test.desc, smells, test.smells)
		}
	}
}
