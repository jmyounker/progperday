package main

import (
	"testing"
)

var parserTests = []struct{
	in string
	want string
}{
	{"-4", "-4"},
	{"4", "4"},
	{"(- 1)", "(- 1)"},
	{"(* 1 2)", "(* 1 2)"},
	{"(/ 1 2)", "(/ 1 2)"},
	{"(+ 1 2)", "(+ 1 2)"},
	{"(- 1 2)", "(- 1 2)"},
	{"(  + 1  2  )", "(+ 1 2)"},
	{"(- (- 3) (+ -14 5))", "(- (- 3) (+ -14 5))"},
	{"(let (x 32) 43)", "(let (x 32) 43)"},
	{"(let (x 32) (* 2 x))", "(let (x 32) (* 2 x))"},
	{"(let (x 32) (let (x x) (* 2 x)))", "(let (x 32) (let (x x) (* 2 x)))"},
}

var errorTests = []string{
	"(* 3)",
	"(1",
	"(let (x x) (* 2 4))",
	"x",
}

func TestEmptyParse(t *testing.T) {
	a, err := parse(newLexer(""))
	if err != nil {
		t.Fatalf("parse error: %s", err)
	}
	if a != nil {
		t.Fatalf("AST not null")
	}
}

func TestParser(t *testing.T) {
	for _, tc := range parserTests {
		a, err := parse(newLexer(tc.in))
		if err != nil {
			t.Fatal("expected no error, but got: %s\n", err)
		}
		if tc.want != a.String() {
			t.Fatalf("want %#v but got %#v", tc.want, a.String())
		}
	}
}

func TestParserErrors(t *testing.T) {
	for _, in := range errorTests {
		_, err := parse(newLexer(in))
		if err == nil {
			t.Fail()
		}
	}
}
