package main

import (
	"fmt"
	"strings"
)

type ast interface {
}

type astExpr interface {
	String() string
}

type astProg struct {
	funcs map[string]*astFnStmt
}

type astNumLit struct {
	value string
}

type astVariable struct {
	value string
}

type astUnaryOpExpr struct {
	arg astExpr
}

type astBinOpExpr struct {
	op   opType
	arg1 astExpr
	arg2 astExpr
}

type astCallExpr struct {
	name string
	args []astExpr
	repr interface{}
}

type astIfElseExpr struct {
	cond astExpr
	ifExpr astExpr
	elseExpr astExpr
}

type astLetExpr struct {
	varName     string
	varInitExpr astExpr
	body        astExpr
}

type astFnStmt struct {
	name string
	args []string
	body astExpr
}

func newNumLitExpr(x string) astExpr {
	return astNumLit{
		value: x,
	}
}

func newVariable(x string) astExpr {
	return astVariable{
		value: x,
	}
}

func newUnaryOpExpr(arg astExpr) astExpr {
	return astCallExpr{
		name: "neg",
		args: []astExpr{arg},
		repr: astUnaryOpExpr{
			arg: arg,
		},
	}
}

func newBinaryOpExpr(fn string, op opType, arg1, arg2 astExpr) astExpr {
	return astCallExpr{
		name: fn,
		args: []astExpr{arg1, arg2},
		repr: astBinOpExpr{
			op: op,
			arg1: arg1,
			arg2: arg2,
		},
	}
}

func newIfElseExpr(cond, ifExpr, elseExpr astExpr) astExpr {
	return astIfElseExpr{
		cond: cond,
		ifExpr: ifExpr,
		elseExpr: elseExpr,
	}
}

func newCallExpr(name string, args []astExpr) *astCallExpr {
	return &astCallExpr{
		name: name,
		args: args,
	}
}

func newLetExpr(name string, initExpr, body astExpr) astExpr {
	return astLetExpr{
		varName:     name,
		varInitExpr: initExpr,
		body:        body,
	}
}

func newFnStmt(name string, args []string, body astExpr) *astFnStmt {
	return &astFnStmt{
		name: name,
		args: args,
		body: body,
	}
}

func (a astProg) String() string {
	r := ""
	for _, fn := range a.funcs {
		r += fn.String()
	}
	return r
}

func (a astFnStmt) String() string {
	return fmt.Sprintf(
		"(fn %s (%s) %s)",
		a.name,
		strings.Join(a.args, " "),
		a.body)
}

func (a astNumLit) String() string {
	return a.value
}

func (a astVariable) String() string {
	return a.value
}

func (a astUnaryOpExpr) String() string {
	return fmt.Sprintf("(- %s)", a.arg)
}

func (a astBinOpExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", a.op, a.arg1, a.arg2)
}

func (a astIfElseExpr) String() string {
	return fmt.Sprintf("(if %s %s %s)", a.cond, a.ifExpr, a.elseExpr)
}

func (a astCallExpr) String() string {
	if len(a.args) == 0 {
		return fmt.Sprintf("(%s)", a.name)
	}
	if a.repr != nil {
		switch r := a.repr.(type) {
		case astUnaryOpExpr:
			return r.String()
		case astBinOpExpr:
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

func (a astLetExpr) String() string {
	return fmt.Sprintf("(let (%s %s) %s)", a.varName, a.varInitExpr, a.body)
}

type opType int

const (
	OP_MINUS opType = iota
	OP_DIV
	OP_PLUS
	OP_MULT
	OP_LT
)

func isKnownOpType(s symbol) bool {
	return s.value == "-" ||
		s.value == "+" ||
		s.value == "*" ||
		s.value == "/" ||
		s.value == "<"
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
	case "<":
		return OP_LT
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
	case OP_LT:
		return "<"
	default:
		panic("must implement translation for op")
	}
}
