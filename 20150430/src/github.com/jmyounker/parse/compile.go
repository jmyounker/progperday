package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/axw/gollvm/llvm"
)

var (
	LLVM_ERR_VAL = llvm.ConstFloat(llvm.DoubleType(), -1.0)
)

func compileAndRun(mod *llvm.Module, f *llvm.Value) float64 {
	llvm.LinkInJIT()
	llvm.InitializeNativeTarget()

	engine, err := llvm.NewJITCompiler(*mod, 2) // optimization level 2
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer engine.Dispose()

	pass := llvm.NewPassManager()
	pass.Add(engine.TargetData())
	pass.AddConstantPropagationPass()
	pass.AddInstructionCombiningPass()
	pass.AddPromoteMemoryToRegisterPass()
	pass.AddGVNPass()
	pass.AddCFGSimplificationPass()
	pass.Run(*mod)

	exec_args := []llvm.GenericValue{}
	exec_res := engine.RunFunction(*f, exec_args)
	return exec_res.Float(llvm.DoubleType())
}

func buildIR(a *astExpr) (*llvm.Module, *llvm.Value, error) {
	mod := llvm.NewModule("calc_module")

	// Define function
	calc_args := []llvm.Type{}
	calc_type := llvm.FunctionType(llvm.DoubleType(), calc_args, false)
	calcf := llvm.AddFunction(mod, "calc", calc_type)
	calcf.SetFunctionCallConv(llvm.CCallConv)

	// Set up basic block for function body
	entry := llvm.AddBasicBlock(calcf, "entry")
	bld := llvm.NewBuilder()
	defer bld.Dispose()
	bld.SetInsertPointAtEnd(entry)

	// Compile expression into basic block
	retv, err := compileExpr(a, "retvn", bld)
	if err != nil {
		return nil, nil, err
	}
	bld.CreateRet(retv)

	// Validate results
	err = llvm.VerifyModule(mod, llvm.ReturnStatusAction)
	if err != nil {
		return nil, nil, fmt.Errorf("error: %s", err)
	}

	return &mod, &calcf, nil
}

func compileExpr(a *astExpr, vn string, bld llvm.Builder) (llvm.Value, error) {
	if a.numLit != nil {
		val, err := strconv.ParseFloat(a.numLit.value, 64)
		if err != nil {
			return LLVM_ERR_VAL, fmt.Errorf("could not convert %s to float", a.numLit.value)
		}
		return llvm.ConstFloat(llvm.DoubleType(), val), nil
	}
	if a.unaryOpExpr != nil {
		expr1, err := compileExpr(a.unaryOpExpr.arg, "iU", bld)
		if err != nil {
			return LLVM_ERR_VAL, err
		}
		return bld.CreateFNeg(expr1, vn), nil
	}
	if a.binOpExpr != nil {
		e := a.binOpExpr
		switch e.op {
		case OP_MINUS:
			return buildBinOp(bld.CreateFSub, e.arg1, e.arg2, vn, bld)
		case OP_PLUS:
			return buildBinOp(bld.CreateFAdd, e.arg1, e.arg2, vn, bld)
		case OP_MULT:
			return buildBinOp(bld.CreateFMul, e.arg1, e.arg2, vn, bld)
		case OP_DIV:
			return buildBinOp(bld.CreateFDiv, e.arg1, e.arg2, vn, bld)
		default:
			panic("an operation has not been implemented yet")
		}
	}
	panic("ast type has not been handled yet")
}


func buildBinOp(f func(llvm.Value, llvm.Value, string) llvm.Value, arg1, arg2 *astExpr, vn string, bld llvm.Builder) (llvm.Value, error) {
	expr1, err := compileExpr(arg1, "iL", bld)
	if err != nil {
		return LLVM_ERR_VAL, err
	}
	expr2, err := compileExpr(arg2, "iR", bld)
	if err != nil {
		return LLVM_ERR_VAL, err
	}
	return f(expr1, expr2, vn), nil

}
