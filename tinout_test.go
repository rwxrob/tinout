package tinout

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	spec, _ := Load("testdata/commonmark-0.29.yaml")
	t.Log(spec)
}

func TestRead(t *testing.T) {
	f, _ := os.Open("testdata/commonmark-0.29.yaml")
	spec, _ := Read(f)
	t.Log(spec)
}

func TestCheck(t *testing.T) {
	alltrue := func(t *Test) bool {
		return t.I == t.I // forces them all
	}
	somefalse := func(t *Test) bool {
		t.Got = "nothing"
		return false
	}
	spec, _ := Load("testdata/commonmark-0.29.yaml")
	result := spec.Check(alltrue)
	if result != nil {
		t.Fail()
	}
	result = spec.Check(somefalse)
	t.Log(result.State())
}
