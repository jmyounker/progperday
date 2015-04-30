package main

// Parser for a very simple list-like calculator language.  Next version
// will probably incorporate a parser, and then subsequently a compiler
// via llvm.

import (
	"errors"
	"fmt"
)

type astExpr struct {
	numLit *astNumLit
	unaryOpExpr *astUnaryOpExpr
	binOpExpr *astBinOpExpr
}

type astNumLit struct {
	value string
}

type astUnaryOpExpr struct {
	arg *astExpr
}

type astBinOpExpr struct {
	op opType
	arg1 *astExpr
	arg2 *astExpr
}

func newNumLitExpr(num string) *astExpr {
	return &astExpr{
		numLit: &astNumLit {
			value: num,
		},
	}
}

func newUnaryOpExpr(arg *astExpr) *astExpr {
	return &astExpr{
		unaryOpExpr: &astUnaryOpExpr {
			arg: arg,
		},
	}
}

func newBinaryOpExpr(op opType, arg1, arg2 *astExpr) *astExpr {
	return &astExpr{
		binOpExpr: &astBinOpExpr {
			op: op,
			arg1: arg1,
			arg2: arg2,
		},
	}
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

func parse(symc chan symbol) (*astExpr, error) {
	p := newParser(symc)
	sym := p.read()
	// horrible, horrible hack
	if sym.symType == SYM_EOF {
		return nil, nil
	}
	return parseExpr(&p)
}
// GRAMMAR
// TODO: implemented "yet" expressions
//
// program := expr
// expr := [ NUM ]
//         [ "(" inside ")" ]
//
// inner := unary | binary | let
// unary := OP("-") expr
// binary := OP expr expr
// let := "let" (SYM expr) expr

func parseExpr(p *parser) (*astExpr, error) {
	sym := p.read()
	switch sym.symType {
	case SYM_LIT:
		p.consume()
		return newNumLitExpr(sym.value), nil
	case SYM_OPEN_PAREN:
		p.consume()
		ast, err := parseInner(p)
		if err != nil {
			return ast, err
		}
		sym := p.read()
		if sym.symType != SYM_CLOSE_PAREN {
			return ast, fmt.Errorf("expected ')' but found %q", sym.value)
		}
		p.consume()
		return ast, err
	default:
		return nil, errors.New("expected '(' or a literal")
	}
}

func parseInner(p *parser) (*astExpr, error) {
	// last sym read was SYM_OPEN_PARENS and it was consumed
	sym := p.read()
	if isKnownOpType(sym) {
		return parseOp(p)
	}
	return nil, errors.New("expected an operation")
}

func parseOp(p *parser) (*astExpr, error) {
	// last sym read was SYM_OP and it is unconsumed
	sym := p.read()
	op := newOpTypeFromSym(sym)
	p.consume()
	args := []*astExpr{}
	for ; sym.symType != SYM_CLOSE_PAREN; sym = p.read() {
		arg, err := parseExpr(p)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	} // terminating parens still unconsumed
	if len(args) == 0 {
		return nil, errors.New("expected at least one argument to operation")
	}
	if len(args) > 2 {
		return nil, errors.New("expected no more than two arguments to an operation")
	}
	if len(args) == 1 {
		if op != OP_MINUS {
			return nil, fmt.Errorf("expected two arguments for operation '%s'", op.String())
		}
		return newUnaryOpExpr(args[0]), nil
	}
	return newBinaryOpExpr(op, args[0], args[1]), nil
}

type parser struct {
	symc chan symbol
	readAhead bool
	next symbol
	astc chan *astExpr
}

func newParser(symc chan symbol) parser {
	return parser{
		symc: symc,
		readAhead: false,
		astc: make(chan *astExpr),
	}
}

func (p *parser)read() symbol {
	if !p.readAhead {
		p.next = <-p.symc
		p.readAhead = true
	}
	return p.next
}

func (p *parser)consume() {
	if p.readAhead {
		p.readAhead = false
	}
}

func (a *astExpr)String() string {
	if a.numLit != nil {
		return a.numLit.String()
	}
	if a.unaryOpExpr != nil {
		return a.unaryOpExpr.String()
	}
	if a.binOpExpr != nil {
		return a.binOpExpr.String()
	}
	panic("malformed ast")
}

func (a *astNumLit)String() string {
	return a.value
}

func (a *astUnaryOpExpr)String() string {
	return fmt.Sprintf("(- %s)", a.arg.String())
}

func (a *astBinOpExpr)String() string {
	return fmt.Sprintf("(%s %s %s)", a.op.String(), a.arg1.String(), a.arg2.String())
}

func (ot opType)String() string {
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
