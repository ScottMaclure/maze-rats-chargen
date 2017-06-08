package main

import "testing"

type Test struct {
	in  int
	out string
}

func Test_contains(t *testing.T) {

	s := []string{"a", "b", "c"}

	result := contains(s, "a")
	if result == false {
		t.Error("Expected true, got false")
	}

	result2 := contains(s, "z")
	if result2 == true {
		t.Error("Expected false, got true")
	}

}

type TestRandInt struct {
	min int
	max int
}

var randIntTests = []TestRandInt{
	{1, 1},
	{1, 3},
	{2, 3},
}

func Test_randInt(t *testing.T) {

	for i, test := range randIntTests {
		r := randInt(test.min, test.max)
		if r < test.min || r > test.max {
			t.Errorf("Test #%d, min=%d, max=%d, result=%d", i, test.min, test.max, r)
		}
	}

}
