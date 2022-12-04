package tinout

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// Spec represents a test specification of inputs and outputs.
type Spec struct {
	Name     string // CommonMark
	Version  string // v0.29
	Source   string // https://gitlab.com/commonmark/commonmark-spec/
	Issues   string // .../issues
	Discuss  string // https://talk.commonmark.org/
	Notes    string // CommonMark is a highly specified Markdown variation.
	Date     string // 2019-04-06
	License  string // http://creativecommons.org/licenses/by-sa/4.0/
	Tests    []Test
	Sections []Section
}

// Section of tests.
type Section struct {
	Name  string
	Notes string
	Tests []Test
}

// Test has the input, output, and notes for a given test.
type Test struct {
	I   string // input
	O   string // output
	N   string // notes
	Got string // result of last check
}

// Passing returns if t.Got is equal to t.O.
func (t *Test) Passing() bool {
	return t.Got == t.O
}

// CheckMethod is any function that takes a test and returns true if the test
// passed. The value it got when testing can be stored in t.Got and t.Passed
// should be set to true if passed.
type CheckMethod func(t *Test) bool

// State returns a string describing the current state of the test.
func (t *Test) State() string {
	passing := "failing"
	if t.Passing() {
		passing = "passing"
	}
	return fmt.Sprintf("\nState:    %q\nInput:    %q\nWanted:   %q\nGot:      %q\n", passing, t.I, t.O, t.Got)
}

// Load loads the Spec from a YAML file at path.
func Load(path string) (Spec, error) {
	s := Spec{}
	byt, err := ioutil.ReadFile(path)
	if err != nil {
		return s, err
	}
	err = yaml.Unmarshal(byt, &s)
	return s, err
}

// Read reads the Spec from a YAML stream reader.
func Read(r io.Reader) (Spec, error) {
	s := Spec{}
	byt, err := ioutil.ReadAll(r)
	if err != nil {
		return s, err
	}
	err = yaml.Unmarshal(byt, &s)
	return s, err
}

// Check takes a function as an argument that takes in CheckMethod function. It
// then calls the CheckMethod on all the s.Tests and s.Sections[].Tests until
// it finds one that does not pass and returns a pointer to it. Returns nil on
// success.
func (s *Spec) Check(ok CheckMethod) *Test {
	for _, t := range s.Tests {
		if !ok(&t) {
			return &t
		}
	}
	for _, sc := range s.Sections {
		for _, t := range sc.Tests {
			if !ok(&t) {
				return &t
			}
		}
	}
	return nil
}
