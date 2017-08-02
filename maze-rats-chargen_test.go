package main

import "testing"

type Test struct {
	in  int
	out string
}

type TestRandInt struct {
	min int
	max int
}

type TestStrings struct {
	in  string
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

var testNLS = []TestStrings{
	{"battle scars", "battle scars"},
	{"birthmark", "a birthmark"},
	{"animal scent", "an animal scent"},
	{"Cautious", "a Cautious"},
	{"  foo  ", "a foo"},
	{"  apple  ", "an apple"},
	{"apples  ", "apples"},
	{"Acid Scars", "Acid Scars"},
	{"scars", "scars"},
}

func Test_naturalLanguageSplice(t *testing.T) {

	initInflections()

	for i, test := range testNLS {
		r := naturalLanguageSplice(test.in)
		if r != test.out {
			t.Errorf("Test #%d, input=%s, expected=%s, actual=%s", i, test.in, test.out, r)
		}
	}

}

func Test_generateSpellName(t *testing.T) {
	r := generateSpellName()
	// TODO implement test
	if r != "foiejef" {
		t.Errorf("Got %s", r)
	}
}
