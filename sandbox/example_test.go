package main

import (
	"errors"
	"testing" // テストで使える関数・構造体が用意されているパッケージをimport
)

// TableDrivenTests ---------------------
type exampleTest struct {
	in       string
	out      int
	errvalue error
}

var exampletests = []exampleTest{
	{"hoge", 1, nil},
	{"fuga", 0, errors.New("code must be hoge")},
}

func TestTabledriven(t *testing.T) {
	for i := range exampletests {
		test := &exampletests[i]
		result, err := example(test.in)

		if test.errvalue == nil { // test for error
			if err != nil {
				t.Fatalf("failed test %#v", err)
			}
		} else {
			if err == nil {
				t.Fatalf("failed test %#v", err)
			}
		}
		if result != test.out { // test for out
			t.Fatal("failed test")
		}
	}
}

// standard tests ------------------------
func TestExampleSuccess(t *testing.T) {
	result, err := example("hoge")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if result != 1 {
		t.Fatal("failed test")
	}
}

func TestExampleFailed(t *testing.T) {
	result, err := example("fuga")
	if err == nil {
		t.Fatal("failed test")
	}
	if result != 0 {
		t.Fatal("failed test")
	}
}
