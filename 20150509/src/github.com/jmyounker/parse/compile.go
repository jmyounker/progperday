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

func compileAndRun(mod *llvm.Module, f llvm.Value) float64 {
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
	exec_res := engine.RunFunction(f, exec_args)
	return exec_res.Float(llvm.DoubleType())
}

type compiler struct {
	bld    llvm.Builder
	varCnt uint64
	symt   symTable
	symft  symTable
	curFunc llvm.Value
}

func newCompiler() *compiler {
	c := compiler{
		bld : llvm.NewBuilder(),
		varCnt: 0,
		symt: symTable{},
		symft: symTable{},
	}
	c.symft.push()

	c.symft.add("neg", funcDefUnary{})
	c.symft.add("plus", funcDefBinary{c.bld.CreateFAdd})
	c.symft.add("minus", funcDefBinary{c.bld.CreateFSub})
	c.symft.add("mult", funcDefBinary{c.bld.CreateFMul})
	c.symft.add("div", funcDefBinary{c.bld.CreateFDiv})
	c.symft.add("lt", funcDefFCmpBinary{llvm.FloatULT})
	return &c
}

type funcDef interface {
	buildCall(c compiler, a astCallExpr, vn string) (llvm.Value, error)
}

type funcDefUnary struct {
}

func (fd funcDefUnary)buildCall(c compiler, a astCallExpr, vn string) (llvm.Value, error) {
	expr, err := c.compileExpr(a.args[0], c.newVar("_i"))
	if err != nil {
		return LLVM_ERR_VAL, err
	}
	return c.bld.CreateFNeg(expr, vn), nil
}

type funcDefBinary struct {
	opcodeBuilder func(llvm.Value, llvm.Value, string) llvm.Value
}

func (fd funcDefBinary)buildCall(c compiler, a astCallExpr, vn string) (llvm.Value, error) {
	expr1, err := c.compileExpr(a.args[0], c.newVar("_iL"))
	if err != nil {
		return LLVM_ERR_VAL, err
	}
	expr2, err := c.compileExpr(a.args[1], c.newVar("_iR"))
	if err != nil {
		return LLVM_ERR_VAL, err
	}
	return fd.opcodeBuilder(expr1, expr2, vn), nil
}

type funcDefFCmpBinary struct {
	cmpPred llvm.FloatPredicate
}

func (fd funcDefFCmpBinary)buildCall(c compiler, a astCallExpr, vn string) (llvm.Value, error) {
	expr1, err := c.compileExpr(a.args[0], c.newVar("_iL"))
	if err != nil {
		return LLVM_ERR_VAL, err
	}
	expr2, err := c.compileExpr(a.args[1], c.newVar("_iR"))
	if err != nil {
		return LLVM_ERR_VAL, err
	}
	return c.bld.CreateFCmp(fd.cmpPred, expr1, expr2, vn), nil
}

type funcDefGeneric struct {
	fn  *astFnStmt
	def llvm.Value
}

func (fd funcDefGeneric)buildCall(c compiler, a astCallExpr, vn string) (llvm.Value, error) {
	args := []llvm.Value{}
	for i, arg := range a.args {
		an := c.newVar(fmt.Sprintf("%s%d", a.name, i))
		res, err := c.compileExpr(arg, an)
		if err != nil {
			return LLVM_ERR_VAL, err
		}
		args = append(args, res)
	}
	return c.bld.CreateCall(fd.def, args, vn), nil
}

func (c compiler) dispose() {
	c.bld.Dispose()
}

func (c compiler) newVar(x string) string {
	c.varCnt += 1
	return fmt.Sprintf("_VAR_x_%d", c.varCnt)
}

func buildIR(a *astProg) (*llvm.Module, llvm.Value, error) {
	mod := llvm.NewModule("calc_module")

	c := newCompiler()
	defer c.dispose()
	for _, fn := range a.funcs {
		// Define function
		fnArgs := []llvm.Type{}
		for i := 0; i < len(fn.args); i++ {
			fnArgs = append(fnArgs, llvm.DoubleType())
		}
		fnType := llvm.FunctionType(llvm.DoubleType(), fnArgs, false)
		fnDef := llvm.AddFunction(mod, fn.name, fnType)
		fnDef.SetFunctionCallConv(llvm.CCallConv)

		c.symft.add(fn.name, &funcDefGeneric{fn: fn, def: fnDef})
	}

	for _, fn := range a.funcs {
		feAny, ok := c.symft.get(fn.name)
		if !ok {
			panic("all functions should have been created here")
		}
		fe := feAny.(*funcDefGeneric)
		c.symt.push()
		c.curFunc = fe.def
		for i := 0; i < len(fe.fn.args); i++ {
			c.symt.add(fe.fn.args[i], fe.def.Param(i))
		}
		entry := llvm.AddBasicBlock(fe.def, fmt.Sprintf("_%s_entry", fn.name))
		c.bld.SetInsertPointAtEnd(entry)

		retv, err := c.compileExpr(fn.body, c.newVar("_retv"))
		if err != nil {
			return nil, llvm.Value{}, err
		}
		c.bld.CreateRet(retv)
		c.curFunc = llvm.Value{} // not defined
		c.symt.pop()
	}

	// Validate results
	if err := llvm.VerifyModule(mod, llvm.ReturnStatusAction); err != nil {
		return nil, llvm.Value{}, fmt.Errorf("error: %s", err)
	}

	feMain, ok := c.symft.get("main")
	if !ok {
		panic("expected to fine a main function")
	}
	return &mod, feMain.(*funcDefGeneric).def, nil
}

func (c compiler) compileExpr(a astExpr, vn string) (llvm.Value, error) {
	switch a := a.(type) {
	case astNumLit:
		val, err := strconv.ParseFloat(a.value, 64)
		if err != nil {
			return LLVM_ERR_VAL, fmt.Errorf("could not convert %s to float", a.value)
		}
		return llvm.ConstFloat(llvm.DoubleType(), val), nil
	case astVariable:
		value, ok := c.symt.get(a.value)
		if !ok {
			return LLVM_ERR_VAL, fmt.Errorf("could not find variable %s", a.value)
		}
		return value.(llvm.Value), nil
	case astCallExpr:
		fe, ok := c.symft.get(a.name)
		if !ok {
			panic(fmt.Sprintf("function %q should already be defined", a.name))
		}
		return fe.(funcDef).buildCall(c, a, vn)
	case astIfElseExpr:
		ifBlock := llvm.AddBasicBlock(c.curFunc, "ifBlock")
		elseBlock := llvm.AddBasicBlock(c.curFunc, "elseBlock")
		mergeBlock := llvm.AddBasicBlock(c.curFunc, "mergeBlock")

		condVal, err := c.compileExpr(a.cond, "condval")
		if err != nil {
			return LLVM_ERR_VAL, err
		}
		c.bld.CreateCondBr(condVal, ifBlock, elseBlock)

		c.bld.SetInsertPointAtEnd(ifBlock)
		ifExitVal, err := c.compileExpr(a.ifExpr, "ifExitVal")
		if err != nil {
			return LLVM_ERR_VAL, err
		}
		c.bld.CreateBr(mergeBlock)
		ifExitBlock := c.bld.GetInsertBlock()

		c.bld.SetInsertPointAtEnd(elseBlock)
		elseExitVal, err := c.compileExpr(a.elseExpr, "elseExitVal")
		if err != nil {
			return LLVM_ERR_VAL, err
		}
		c.bld.CreateBr(mergeBlock)
		elseExitBlock := c.bld.GetInsertBlock()

		c.bld.SetInsertPointAtEnd(mergeBlock)
		inVals := []llvm.Value{ifExitVal, elseExitVal}
		inBlocks := []llvm.BasicBlock{ifExitBlock, elseExitBlock}
		phi := c.bld.CreatePHI(llvm.DoubleType(), vn)
		phi.AddIncoming(inVals, inBlocks)
		return phi, nil
	case astLetExpr:
		v := c.newVar(a.varName)
		val, err := c.compileExpr(a.varInitExpr, v)
		if err != nil {
			return val, err
		}
		c.symt.push()
		c.symt.add(a.varName, val)
		defer c.symt.pop()
		return c.compileExpr(a.body, vn)
	default:
		panic(fmt.Sprintf("ast type %#v has not been handled yet", a))
	}
}
