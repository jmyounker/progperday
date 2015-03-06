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

func buildIR(a *ast) (*llvm.Module, *llvm.Value, error) {
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

func compileExpr(a *ast, vn string, bld llvm.Builder) (llvm.Value, error) {
	switch (a.astType) {
	case AST_LIT:
		// Move value conversion to parsing stage
		val, err := strconv.ParseFloat(a.value.value, 64)
		if err != nil {
			return LLVM_ERR_VAL, fmt.Errorf("could not convert %s to float", a.value.value)
		}
		return llvm.ConstFloat(llvm.DoubleType(), val), nil

	case AST_NEG:
		expr1, err := compileExpr(a.op1, "iU", bld)
		if err != nil {
			return LLVM_ERR_VAL, err
		}
		return bld.CreateFNeg(expr1, vn), nil
	case AST_MULT:
		return buildBinOp(bld.CreateFMul, a.op1, a.op2, vn, bld)
	case AST_DIV:
		return buildBinOp(bld.CreateFDiv, a.op1, a.op2, vn, bld)
	case AST_PLUS:
		return buildBinOp(bld.CreateFAdd, a.op1, a.op2, vn, bld)
	case AST_MINUS:
		return buildBinOp(bld.CreateFSub, a.op1, a.op2, vn, bld)
	default:
		panic("unknown AST type encountered")
	}
}

func buildBinOp(f func(llvm.Value, llvm.Value, string) llvm.Value, op1, op2 *ast, vn string, bld llvm.Builder) (llvm.Value, error) {
	expr1, err := compileExpr(op1, "iL", bld)
	if err != nil {
		return LLVM_ERR_VAL, err
	}
	expr2, err := compileExpr(op2, "iR", bld)
	if err != nil {
		return LLVM_ERR_VAL, err
	}
	return f(expr1, expr2, vn), nil

}
