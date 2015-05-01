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
	variable *astVariable
	unaryOpExpr *astUnaryOpExpr
	binOpExpr *astBinOpExpr
	letExpr *astLetExpr
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
	op opType
	arg1 *astExpr
	arg2 *astExpr
}

type astLetExpr struct {
	varName string
	varInitExpr *astExpr
	body *astExpr
}

func newNumLitExpr(x string) *astExpr {
	return &astExpr{
		numLit: &astNumLit {
			value: x,
		},
	}
}

func newVariable(x string) *astExpr {
	return &astExpr{
		variable: &astVariable {
			value: x,
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

func newLetExpr(name string, initExpr, body *astExpr) *astExpr {
	return &astExpr{
		letExpr: &astLetExpr{
			varName: name,
			varInitExpr: initExpr,
			body: body,
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
//
// program := expr
// expr := [ NUM ]
//         [ SYM ]
//         [ "(" inside ")" ]
//
// inner := unary | binary | let
// unary := OP("-") expr
// binary := OP expr expr
// let := "let" (SYM expr) expr

func parseExpr(p *parser) (*astExpr, error) {
	// terminates successfully after consuming all tokens in the expression
	sym := p.read()
	switch sym.symType {
	case SYM_LIT:
		p.consume()
		return newNumLitExpr(sym.value), nil
	case SYM_SYM:
		if _, ok := p.symt.get(sym.value); !ok {
			return nil, fmt.Errorf("variable %s is not defined", sym.value)
		}
		p.consume()
		return newVariable(sym.value), nil
	case SYM_OPEN_PAREN:
		p.consume()
		ast, err := parseInner(p)
		if err != nil {
			return ast, err
		}
		sym := p.read()
		if sym.symType != SYM_CLOSE_PAREN {
			return ast, fmt.Errorf("expected ')' but found '%s'", sym.value)
		}
		p.consume()
		return ast, err
	default:
		return nil, errors.New("expected '(' or a literal")
	}
}

func parseInner(p *parser) (*astExpr, error) {
	// last sym read was SYM_OPEN_PARENS and it was consumed
	// terminates with final ')' unconsumed
	sym := p.read()
	if isKnownOpType(sym) {
		return parseOp(p)
	}
	if sym.symType == SYM_SYM && sym.value == "let" {
		return parseLetExpr(p)
	}
	return nil, errors.New("expected an operation or 'let'")
}

func parseOp(p *parser) (*astExpr, error) {
	// last sym read was SYM_OP and it is unconsumed
	// terminates with final ')' unconsumed
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

func parseLetExpr(p *parser) (*astExpr, error) {
	// last sym read was SYM_LET and it is unconsumed
	// terminates successfully with final ')' still unconsumed
	p.consume()  // eat 'let'
	sym := p.read()
	if sym.symType != SYM_OPEN_PAREN {
		return nil, errors.New("expected '(' at beginning of variable clause")
	}
	p.consume() // eat '('
	sym = p.read()
	if sym.symType != SYM_SYM {
		return nil, fmt.Errorf("expected a variable name and not %q", sym)
	}
	varName := sym.value
	p.consume() // eat variable name
	varInitExpr, err := parseExpr(p)
	if err != nil {
		 return nil, err
	}
	// The symbol table entry must be made *after* the init parsing, but before the body
	// parsing.
	p.symt.push()
	p.symt.add(varName, struct{}{})
	defer p.symt.pop()
	sym = p.read()
	if sym.symType != SYM_CLOSE_PAREN {
		return nil, fmt.Errorf("expected variable definition to end with ')' and not %q", sym)
	}
	p.consume()
	body, err := parseExpr(p)
	if err != nil {
		return nil, err
	}
	return newLetExpr(varName, varInitExpr, body), nil
}

type parser struct {
	symc chan symbol
	readAhead bool
	next symbol
	astc chan *astExpr
	symt symTable
}

func newParser(symc chan symbol) parser {
	return parser{
		symc: symc,
		readAhead: false,
		astc: make(chan *astExpr),
		symt: symTable{},
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
	if a == nil {
		return "NIL"
	}
	if a.numLit != nil {
		return a.numLit.String()
	}
	if a.variable != nil {
		return a.variable.String()
	}
	if a.unaryOpExpr != nil {
		return a.unaryOpExpr.String()
	}
	if a.binOpExpr != nil {
		return a.binOpExpr.String()
	}
	if a.letExpr != nil {
		return a.letExpr.String()
	}
	panic("malformed ast")
}

func (a *astNumLit)String() string {
	return a.value
}

func (a *astVariable)String() string {
	return a.value
}

func (a *astUnaryOpExpr)String() string {
	return fmt.Sprintf("(- %s)", a.arg)
}

func (a *astBinOpExpr)String() string {
	return fmt.Sprintf("(%s %s %s)", a.op, a.arg1, a.arg2)
}

func (a *astLetExpr)String() string {
	return fmt.Sprintf("(let (%s %s) %s)", a.varName, a.varInitExpr, a.body)
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
