package main

import (
	"fmt"
	"strings"
)

type astProg struct {
	funcs map[string]*astFnStmt
}

type astExpr struct {
	numLit      *astNumLit
	variable    *astVariable
	callExpr    *astCallExpr
	letExpr     *astLetExpr
}

type astNumLit struct {
	value string
}

type astVariable struct {
	value string
}

type astUnaryOpExpr struct {
	arg *astExpr
}

type astBinOpExpr struct {
	op   opType
	arg1 *astExpr
	arg2 *astExpr
}

type astCallExpr struct {
	name string
	args []*astExpr
	repr interface{}
}

type astLetExpr struct {
	varName     string
	varInitExpr *astExpr
	body        *astExpr
}

type astFnStmt struct {
	name string
	args []string
	body *astExpr
}

func newNumLitExpr(x string) *astExpr {
	return &astExpr{
		numLit: &astNumLit{
			value: x,
		},
	}
}

func newVariable(x string) *astExpr {
	return &astExpr{
		variable: &astVariable{
			value: x,
		},
	}
}

func newUnaryOpExpr(arg *astExpr) *astExpr {
	return &astExpr{
		callExpr: &astCallExpr{
			name: "neg",
			args: []*astExpr{arg},
			repr: &astUnaryOpExpr{
				arg: arg,
			},
		},
	}
}

func newBinaryOpExpr(fn string, op opType, arg1, arg2 *astExpr) *astExpr {
	return &astExpr{
		callExpr: &astCallExpr{
			name: fn,
			args: []*astExpr{arg1, arg2},
			repr: &astBinOpExpr{
				op: op,
				arg1: arg1,
				arg2: arg2,
			},
		},
	}
}

func newCallExpr(name string, args []*astExpr) *astExpr {
	return &astExpr{
		callExpr: &astCallExpr{
			name: name,
			args: args,
		},
	}
}

func newLetExpr(name string, initExpr, body *astExpr) *astExpr {
	return &astExpr{
		letExpr: &astLetExpr{
			varName:     name,
			varInitExpr: initExpr,
			body:        body,
		},
	}
}

func newFnStmt(name string, args []string, body *astExpr) *astFnStmt {
	return &astFnStmt{
		name: name,
		args: args,
		body: body,
	}
}

func (a *astProg) String() string {
	if a == nil {
		return "NIL"
	}
	r := ""
	for _, fn := range a.funcs {
		r += fn.String()
	}
	return r
}

func (a *astFnStmt) String() string {
	return fmt.Sprintf(
		"(fn %s (%s) %s)",
		a.name,
		strings.Join(a.args, " "),
		a.body)
}

func (a *astExpr) String() string {
	if a.numLit != nil {
		return a.numLit.String()
	}
	if a.variable != nil {
		return a.variable.String()
	}
	if a.callExpr != nil {
		return a.callExpr.String()
	}
	if a.letExpr != nil {
		return a.letExpr.String()
	}
	panic("malformed ast")
}

func (a *astNumLit) String() string {
	return a.value
}

func (a *astVariable) String() string {
	return a.value
}

func (a *astUnaryOpExpr) String() string {
	return fmt.Sprintf("(- %s)", a.arg)
}

func (a *astBinOpExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", a.op, a.arg1, a.arg2)
}

func (a *astCallExpr) String() string {
	if len(a.args) == 0 {
		return fmt.Sprintf("(%s)", a.name)
	}
	if a.repr != nil {
		switch r := a.repr.(type) {
		case *astUnaryOpExpr:
			return r.String()
		case *astBinOpExpr:
			return r.String()
		default:
			panic("unknown expression simpification")
		}
	}
	args := []string{}
	for _, arg := range a.args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("(%s %s)", a.name, strings.Join(args, " "))
}

func (a *astLetExpr) String() string {
	return fmt.Sprintf("(let (%s %s) %s)", a.varName, a.varInitExpr, a.body)
}

type opType int

const (
	OP_MINUS opType = iota
	OP_DIV
	OP_PLUS
	OP_MULT
)

func isKnownOpType(s symbol) bool {
	return s.value == "-" || s.value == "+" || s.value == "*" || s.value == "/"
}

func newOpTypeFromSym(s symbol) opType {
	switch s.value {
	case "-":
		return OP_MINUS
	case "+":
		return OP_PLUS
	case "*":
		return OP_MULT
	case "/":
		return OP_DIV
	default:
		panic("not reachable")
	}
}

func (ot opType) String() string {
	switch ot {
	case OP_MINUS:
		return "-"
	case OP_PLUS:
		return "+"
	case OP_MULT:
		return "*"
	case OP_DIV:
		return "/"
	default:
		panic("must implement translation for op")
	}
}
