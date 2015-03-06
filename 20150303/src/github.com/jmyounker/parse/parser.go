package main

// Parser for a very simple list-like calculator language.  Next version
// will probably incorporate a parser, and then subsequently a compiler
// via llvm.

import (
	"errors"
	"fmt"
)

type ast struct {
	astType astType
	op1 *ast
	op2 *ast
	value symbol
}

type astType int

const (
	AST_MULT astType = iota
	AST_DIV
	AST_PLUS
	AST_MINUS
	AST_NEG
	AST_LIT
)

func parse(symc chan symbol) (*ast, error) {
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
func parseExpr(p *parser) (*ast, error) {
	sym := p.read()
	p.consume()
	if sym.symType == SYM_LIT {
		return &ast{
			astType: AST_LIT,
			value: sym,
		}, nil
	}
	if sym.symType == SYM_OPEN_PAREN {
		sym = p.read()
		p.consume()
		op := sym.value
		if sym.symType == SYM_OP && op == "-" {
			op1, err := parseExpr(p)
			if err != nil {
				return nil, err
			}
			sym = p.read()
			if sym.symType == SYM_CLOSE_PAREN {
				p.consume()
				return &ast{
					astType: AST_NEG,
					op1: op1,
				}, nil
			}
			op2, err := parseExpr(p)
			if err != nil {
				return nil, err
			}
			sym = p.read()
			p.consume()
			if sym.symType != SYM_CLOSE_PAREN {
				return nil, errors.New("expected close paren")
			}
			return &ast{
				astType: AST_MINUS,
				op1: op1,
				op2: op2,
			}, nil
		}
		op1, err := parseExpr(p)
		if err != nil {
			return nil, err
		}
		op2, err := parseExpr(p)
		if err != nil {
			return nil, err
		}
		sym = p.read()
		p.consume()
		if sym.symType != SYM_CLOSE_PAREN {
			return nil, errors.New("expected close paren")
		}
		if op == "*" {
			return &ast{
				astType: AST_MULT,
				op1: op1,
				op2: op2,
			}, nil
		}
		if op == "/" {
			return &ast{
				astType: AST_DIV,
				op1: op1,
				op2: op2,
			}, nil
		}
		if op == "+" {
			return &ast{
				astType: AST_PLUS,
				op1: op1,
				op2: op2,
			}, nil
		}
		return nil, errors.New("unknown symbol")
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

func (a *ast)String() string {
	switch a.astType {
	case AST_LIT:
		return a.value.value
	case AST_NEG:
		return fmt.Sprintf("(- %s)", a.op1.String())
	case AST_MULT:
		return fmt.Sprintf("(* %s %s)", a.op1.String(), a.op2.String())
	case AST_DIV:
		return fmt.Sprintf("(/ %s %s)", a.op1.String(), a.op2.String())
	case AST_PLUS:
		return fmt.Sprintf("(+ %s %s)", a.op1.String(), a.op2.String())
	case AST_MINUS:
		return fmt.Sprintf("(- %s %s)", a.op1.String(), a.op2.String())
	default:
		panic("unknown decoding")
	}
}
