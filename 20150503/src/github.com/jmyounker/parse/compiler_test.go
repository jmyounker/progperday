package main

import (
	"testing"
	"fmt"
)

var compilerTests = []struct {
	in   string
	want float64
}{
//	{"(fn main () -4)", -4.0},
//	{"4", 4},
//	{"(- 1)", -1},
//	{"(* 1 2)", 2.0},
//	{"(/ 1 2)", 0.5},
//	{"(+ 1 2)", 3.0},
//	{"(- 1 2)", -1.0},
//	{"(- (- 3) (+ -14 5))", 6.0},
//	{"(let (x 32) 43)", 43.0},
//	{"(let (x 32) x)", 32.0},
//	{"(let (x 32) (* 2 x))", 64.0},
//	{"(let (x 32) (let (x x) (* 2 x)))", 64.0},
//	{"(let (x 32) (+ (let (x 64) x) x))", 96.0},
//	{"(fn main () (let (x 32) (let (y (/ x 16)) (* x y))))", 64.0},
	{"(fn main () (add 2 4)) (fn add (x y) (+ x y))", 6.0},

}

func TestCompile(t *testing.T) {
	for _, tc := range compilerTests {
		a, err := parse(newLexer(tc.in))
		if err != nil {
			t.Fatal("expected no error, but got: %s\n", err)
		}
		mod, f, err := buildIR(a)
		if err != nil {
			fmt.Printf("%s\n", a)
			t.Fatalf("failed to compile code to IR:", err)
		}
		res := compileAndRun(mod, f)
		if tc.want != res {
			t.Fatalf("want %#v but got %#v", tc.want, res)
		}
	}
}
