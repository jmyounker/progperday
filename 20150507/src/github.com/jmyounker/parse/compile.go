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
}

func newCompiler() *compiler {
	c := compiler{
		bld : llvm.NewBuilder(),
		varCnt: 0,
		symt: symTable{},
		symft: symTable{},
	}
	c.symft.push()
	return &c
}

type funcEntry struct {
	fn  *astFnStmt
	def llvm.Value
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

		c.symft.add(fn.name, &funcEntry{fn: fn, def: fnDef})
	}

	for _, fn := range a.funcs {
		feAny, ok := c.symft.get(fn.name)
		if !ok {
			panic("all functions should have been created here")
		}
		fe := feAny.(*funcEntry)
		c.symt.push()
		for i := 0; i < len(fe.fn.args); i++ {
			c.symt.add(fe.fn.args[i], fe.def.Param(i))
		}
		entry := llvm.AddBasicBlock(fe.def, fmt.Sprintf("_%s_entry", fe.fn.name))
		c.bld.SetInsertPointAtEnd(entry)

		retv, err := c.compileExpr(fn.body, c.newVar("_retv"))
		if err != nil {
			return nil, llvm.Value{}, err
		}
		c.bld.CreateRet(retv)
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
	return &mod, feMain.(*funcEntry).def, nil
}

func (c compiler) compileProg(a *astProg) error {
	for _, fn := range a.funcs {
		if err := c.compileFunc(fn); err != nil {
			return err
		}
	}
	return nil
}

func (c compiler) compileFunc(a *astFnStmt) error {
 	return nil
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
	if a.callExpr != nil {
		fe, ok := c.symft.get(a.callExpr.name)
		if !ok {
			panic(fmt.Sprintf("function %q should already be defined", a.callExpr.name))
		}
		args := []llvm.Value{}
		for i, arg := range a.callExpr.args {
			an := c.newVar(fmt.Sprintf("%s%d", a.callExpr.name, i))
			res, err := c.compileExpr(arg, an)
			if err != nil {
				return LLVM_ERR_VAL, err
			}
			args = append(args, res)
		}
		def := fe.(*funcEntry).def
		return c.bld.CreateCall(def, args, vn), nil
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
	expr1, err := c.compileExpr(arg1, c.newVar("_iL"))
	if err != nil {
		return LLVM_ERR_VAL, err
	}
	expr2, err := c.compileExpr(arg2, c.newVar("_iR"))
	if err != nil {
		return LLVM_ERR_VAL, err
	}
	return f(expr1, expr2, vn), nil
}
