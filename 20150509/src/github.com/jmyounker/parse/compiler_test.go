package main

import (
	"testing"
	"fmt"
	"github.com/axw/gollvm/llvm"
)

var retval = llvm.ConstFloat(llvm.DoubleType(), 22.0)

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
	{"(fn add (x y) (+ x y)) (fn main () (add 2 4))", 6.0},
	{"(fn main () (f)) (fn f () 34)", 34.0},
	{"(fn main () (if (< 1 20) -17 42))", -17},
	{"(fn main () (if (< 20 1) -17 42))", 42},
	{"(fn main () (if (< 1 20) (if (< 1 20) -17 23) 42))", -17},
	{"(fn main () (if (< 1 20) -17 (if (< 1 20) 23 42)))", -17},
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


// Skeleton for experiments
//
//func TestBuildFunc(t *testing.T) {
//	mod := llvm.NewModule("calc_module")
//
//	fnArgs := []llvm.Type{llvm.DoubleType(), llvm.DoubleType()}
//	fnType := llvm.FunctionType(llvm.DoubleType(), fnArgs, false)
//	fnDef := llvm.AddFunction(mod, "negate", fnType)
//	fnDef.SetFunctionCallConv(llvm.CCallConv)
//
//	// (fn main () (add 2 3))
//	mainArgs := []llvm.Type{}
//	mainType := llvm.FunctionType(llvm.DoubleType(), mainArgs, false)
//	mainDef := llvm.AddFunction(mod, "main", mainType)
//	mainDef.SetFunctionCallConv(llvm.CCallConv)
//	mainEntry := llvm.AddBasicBlock(mainDef, "main_entry")
//	bld := llvm.NewBuilder()
//	bld.SetInsertPointAtEnd(mainEntry)
//	v1 := llvm.ConstFloat(llvm.DoubleType(), 1.0)
//	v2 := llvm.ConstFloat(llvm.DoubleType(), 2.0)
//	callArgs := []llvm.Value{v1, v2}
//	res := bld.CreateCall(fnDef, callArgs, "retval")
//	bld.CreateRet(res)
//
// 	// (fn add (x y) (+ x y))
//	fnDef.SetFunctionCallConv(llvm.CCallConv)
//	_ = llvm.AddBasicBlock(fnDef, "negate_entry")
//	bld.SetInsertPointAtEnd(mainEntry)
//	fmt.Printf("p0: %s\n", fnDef.Param(0))
//	fmt.Printf("p1: %s\n", fnDef.Param(1))
//	retval := bld.CreateFAdd(fnDef.Param(0), fnDef.Param(1), "retval")
//	bld.CreateRet(retval)
//}
