package golang

import (
	"os"
	"sort"
	"testing"
)

var codeBad = `package bad

func ok() {
	println(3 % 2 == 0)
}
`

var codeGood = `package good

// an instance in time
func ok() {
	println(3%2 == 0)
}
`

var codeNewlineBefore = `
package nl

func ok() {
	println(3%2 == 0)
}
`

var codeNewlineAfter = `package nl

func ok() {
	println(3%2 == 0)
}`

var codeWindowsNewline = "package win\r\n\r\nimport \"time\"\r\n\r\nconst a = 1"

var codeUnreachable = `package vet

func ok() bool {
	return true
	return false
}
`

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

var codeBuild = `// +build !example

package bc

func ok() bool {
	return true
}
`

var codeSnake = `package snake

func snake_case() bool {
	return true
}
`

var codeMixed = `package mixed

func mixedCaps() bool {
	return true
}
`

var codeScream = `package scream

const SCREAMING_SNAKE = "scream"
`

var codeZero = `package zero

var i int = 0
`

var codeOutdent = `package outdent

func ok(i int) bool {
	if i == 0 {
		return true
	} else {
		return false
	}
}
`

var codeInstanceBad = `package instance

type thing int

// New creates an instance of thing
func New() thing {
	return thing(0)
}
`

var codeObjectBad = `package instance

type thing int

// New generates a new thing object
func New() thing {
	return thing(0)
}
`

var codeReceiverName = `
package consistent

type Score int

func (s1 Score) Incr() Score {
	return s + 1
}

func (s2 Score) Decr() Score {
	return s - 1
}
`

var codeRangeSlice = `
package loop

func do() {
	for i, _ := range []int{1, 1, 2, 3, 5} {
		// do something
	}
}
`

var codeRangeMap = `
package loop

func do() {
	for k, _ := range map[string]int{"alice": 9, "bob": 12} {
		// do something
	}
}
`

var codeBadComment = `
package doc

// It's all good. Or perhaps not.
func OK() bool {
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
		{"top", codeNewlineBefore, true},
		{"bottom", codeNewlineAfter, true},
		{"win", codeWindowsNewline, true},
	}

	for _, test := range tests {
		s := newSolution(map[string]string{"code.go": test.code})
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

func TestVetted(t *testing.T) {
	var tests = []struct {
		desc, code string
		ok         bool
	}{
		{"good", codeGood, true},
		{"unreachable", codeUnreachable, false},
	}

	for _, test := range tests {
		s := newSolution(map[string]string{"code.go": test.code})
		if err := s.write(); err != nil {
			t.Fatal(err)
		}
		defer os.Remove(s.dir)

		ok, err := isVetted(s)
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
		s := newSolution(map[string]string{test.desc + `.go`: test.code})

		if err := s.extractComments(); err != nil {
			t.Fatal(err)
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

func TestBuildDirective(t *testing.T) {
	var tests = []struct {
		desc, code string
		ok         bool
	}{
		{"good", codeGood, true},
		{"build", codeBuild, false},
	}
	for _, test := range tests {
		s := newSolution(map[string]string{test.desc + `.go`: test.code})

		if err := s.extractComments(); err != nil {
			t.Fatal(err)
		}

		ok, err := noBuildConstraint(s)
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
		{"build", codeBuild, []string{"build-constraint"}},
		{"snake", codeSnake, []string{"mixed-caps"}},
		{"scream", codeScream, []string{"mixed-caps"}},
		{"unreachable", codeUnreachable, []string{"go-vet"}},
		{"zero", codeZero, []string{"zero-value"}},
		{"outdent", codeOutdent, []string{"if-return-else"}},
		{"instance", codeInstanceBad, []string{"instance"}},
		{"object", codeObjectBad, []string{"object"}},
		{"receiver name", codeReceiverName, []string{"receiver-name"}},
		{"range with map", codeRangeMap, []string{"range-loop"}},
		{"range with slice", codeRangeSlice, []string{"range-loop"}},
		{"incorrect comment format", codeBadComment, []string{"comment-format"}},
	}

	for _, test := range tests {
		smells, err := Analyze(map[string]string{"code.go": test.code})
		if err != nil {
			t.Fatal(err)
		}
		if len(test.smells) != len(smells) {
			t.Errorf("%s: got %v, want %v", test.desc, smells, test.smells)
			continue
		}

		sort.Strings(smells)
		sort.Strings(test.smells)
		for i := 0; i < len(test.smells); i++ {
			if smells[i] != test.smells[i] {
				t.Errorf("%s: got %s, want %v", test.desc, smells[i], test.smells[i])
			}
		}
	}
}
