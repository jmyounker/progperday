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
	p := &parser{}
	p.symc = symc
	p.readAhead = false
	sym := p.read()
	// horrible, horrible hack
	if sym.symType == SYM_EOF {
		return nil, nil
	}
	return parseExpr(p)
}

// Quick and dirty parser.  Were the language any more complicated then I'd
// probably code up a state machine.
func parseExpr(p *parser) (*astExpr, error) {
	sym := p.read()
	p.consume()
	if sym.symType == SYM_LIT {
		return newNumLitExpr(sym.value), nil
	}
	if sym.symType == SYM_OPEN_PAREN {
		sym := p.read()
		if !isKnownOpType(sym) {
			return nil, errors.New("unexpected symbol")
		}
		p.consume()
		op := newOpTypeFromSym(sym)
		if sym.symType == SYM_OP && op == OP_MINUS {
			arg1, err := parseExpr(p)
			if err != nil {
				return nil, err
			}
			sym = p.read()
			if sym.symType == SYM_CLOSE_PAREN {
				p.consume()
				return newUnaryOpExpr(arg1), nil
			}

			arg2, err := parseExpr(p)
			if err != nil {
				return nil, err
			}
			sym = p.read()
			p.consume()
			if sym.symType != SYM_CLOSE_PAREN {
				return nil, errors.New("expected close paren")
			}
			return newBinaryOpExpr(op, arg1, arg2), nil
		}
		arg1, err := parseExpr(p)
		if err != nil {
			return nil, err
		}
		arg2, err := parseExpr(p)
		if err != nil {
			return nil, err
		}
		sym = p.read()
		p.consume()
		if sym.symType != SYM_CLOSE_PAREN {
			return nil, errors.New("expected close paren")
		}
		return newBinaryOpExpr(op, arg1, arg2), nil
	}
	return nil, errors.New("expected '(' or a literal")
}

type parser struct {
	symc chan symbol
	readAhead bool
	next symbol
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
