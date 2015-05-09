package main

import (
	"testing"
)

var exprTests = []struct{
	in string
	want string
}{
	{"-4", "-4"},
	{"4", "4"},
	{"(- 1)", "(- 1)"},
	{"(-\t1\n)", "(- 1)"},
	{"(* 1 2)", "(* 1 2)"},
	{"(/ 1 2)", "(/ 1 2)"},
	{"(+ 1 2)", "(+ 1 2)"},
	{"(- 1 2)", "(- 1 2)"},
	{"(  + 1  2  )", "(+ 1 2)"},
	{"(- (- 3) (+ -14 5))", "(- (- 3) (+ -14 5))"},
	{"(let (x 32) 43)", "(let (x 32) 43)"},
	{"(let (x 32) (* 2 x))", "(let (x 32) (* 2 x))"},
	{"(let (x 32) (let (x x) (* 2 x)))", "(let (x 32) (let (x x) (* 2 x)))"},
	{"(foo 43 47 36)", "(foo 43 47 36)"},
	{"(+ (bar) 37)", "(+ (bar) 37)"},
	{"(< 37 23)", "(< 37 23)"},
	{"(if (< 26 32) 0 1)", "(if (< 26 32) 0 1)"},
}


var progTests = []struct{
	in string
	want string
}{
	{"(fn add (x y) (+ x y)) (fn main () (add 2 4)) ", "(fn add (x y) (+ x y))(fn main () (add 2 4))"},
}

var errorTests = []string{
	"(* 3)",
	"(1",
	"(let (x x) (* 2 4))",
	"x",
	"(+ (let (x 4) (* 2 4)) x)",
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

func TestExprParser(t *testing.T) {
	for _, tc := range exprTests {
		a, err := parseExprFrom(newLexer(tc.in))
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
		_, err := parseExprFrom(newLexer(in))
		if err == nil {
			t.Fail()
		}
	}
}

func TestProgParser(t *testing.T) {
	for _, tc := range progTests {
		a, err := parse(newLexer(tc.in))
		if err != nil {
			t.Fatal("expected no error, but got: %s\n", err)
		}
		if tc.want != a.String() {
			t.Fatalf("want %#v but got %#v", tc.want, a.String())
		}
	}
}

func parseExprFrom(symc chan symbol) (astExpr, error) {
	p := newParser(symc)
	sym := p.read()
	// horrible, horrible hack
	if sym.symType == SYM_EOF {
		return nil, nil
	}
	return parseExpr(&p)
}
