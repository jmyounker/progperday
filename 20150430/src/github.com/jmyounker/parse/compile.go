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

type compiler struct {
	bld    llvm.Builder
	varCnt uint64
	symt   symTable
}

func (c compiler) dispose() {
	c.bld.Dispose()
}

func (c compiler) newVar(x string) string {
	c.varCnt += 1
	return fmt.Sprintf("_VAR_x_%d", c.varCnt)
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
	symt := symTable{}
	c := compiler{
		bld:  bld,
		symt: symt,
	}
	defer c.dispose()
	bld.SetInsertPointAtEnd(entry)

	// Compile expression into basic block
	retv, err := c.compileExpr(a, "_retvn")
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

func (c compiler) compileExpr(a *astExpr, vn string) (llvm.Value, error) {
	if a.numLit != nil {
		val, err := strconv.ParseFloat(a.numLit.value, 64)
		if err != nil {
			return LLVM_ERR_VAL, fmt.Errorf("could not convert %s to float", a.numLit.value)
		}
		return llvm.ConstFloat(llvm.DoubleType(), val), nil
	}
	if a.variable != nil {
		value, ok := c.symt.get(a.variable.value)
		if !ok {
			return LLVM_ERR_VAL, fmt.Errorf("could not find variable %s", a.variable.value)
		}
		return value.(llvm.Value), nil
	}
	if a.unaryOpExpr != nil {
		expr1, err := c.compileExpr(a.unaryOpExpr.arg, "iU")
		if err != nil {
			return LLVM_ERR_VAL, err
		}
		return c.bld.CreateFNeg(expr1, vn), nil
	}
	if a.binOpExpr != nil {
		e := a.binOpExpr
		switch e.op {
		case OP_MINUS:
			return c.buildBinOp(c.bld.CreateFSub, e.arg1, e.arg2, vn)
		case OP_PLUS:
			return c.buildBinOp(c.bld.CreateFAdd, e.arg1, e.arg2, vn)
		case OP_MULT:
			return c.buildBinOp(c.bld.CreateFMul, e.arg1, e.arg2, vn)
		case OP_DIV:
			return c.buildBinOp(c.bld.CreateFDiv, e.arg1, e.arg2, vn)
		default:
			panic("an operation has not been implemented yet")
		}
	}
	if a.letExpr != nil {
		v := c.newVar(a.letExpr.varName)
		val, err := c.compileExpr(a.letExpr.varInitExpr, v)
		if err != nil {
			return val, err
		}
		c.symt.push()
		c.symt.add(a.letExpr.varName, val)
		defer c.symt.pop()
		return c.compileExpr(a.letExpr.body, vn)
	}
	panic("ast type has not been handled yet")
}

func (c compiler) buildBinOp(f func(llvm.Value, llvm.Value, string) llvm.Value, arg1, arg2 *astExpr, vn string) (llvm.Value, error) {
	expr1, err := c.compileExpr(arg1, "iL")
	if err != nil {
		return LLVM_ERR_VAL, err
	}
	expr2, err := c.compileExpr(arg2, "iR")
	if err != nil {
		return LLVM_ERR_VAL, err
	}
	return f(expr1, expr2, vn), nil
}
