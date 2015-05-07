package main

// Parser for a very for a slowly growing language.

import (
	"errors"
	"fmt"
)

func parse(symc chan symbol) (*astProg, error) {
	p := newParser(symc)
	sym := p.read()
	// horrible, horrible hack
	if sym.symType == SYM_EOF {
		return nil, nil
	}
	return parseProg(&p)
}

// GRAMMAR
//
// program := fn_def+
// fn_def := "(" "fn" SYM "(" SYM* ")" expr ")"

// expr := [ NUM ]
//         [ SYM ]
//         [ "(" inside ")" ]
//
// inner := unary | binary | call | let
// unary := OP("-") expr
// binary := OP expr expr
// call := SYM expr*
// let := "let" "(" SYM expr ")" expr

func parseProg(p *parser) (*astProg, error) {
	parseLoop:
	for true {
		sym := p.read()

		switch sym.symType {
		case SYM_EOF:
			break parseLoop

		case SYM_OPEN_PAREN:
			fn, err := parseFnStmt(p)
			if err != nil {
				return nil, err
			}

			// record function definition
			if _, ok := p.funcs[fn.name]; ok {
				return nil, fmt.Errorf("function %s is already defined", fn.name)
			}
			p.funcs[fn.name] = fn

		default:
			return nil, fmt.Errorf("expected '(' to start function definition, but got %s", sym.value)
		}
	}

	// ensure that main function exists
	fnMain, ok := p.funcs["main"]
	if !ok {
		return nil, errors.New("a function 'main' must be defined")
	}
	if len(fnMain.args) != 0 {
		return nil, fmt.Errorf(
			"function 'main' has %d arguments but must have zero",
			len(fnMain.args))
	}

	if err := resolveFuncs(p); err != nil {
		return nil, err
	}

	return &astProg{funcs: p.funcs}, nil
}

func parseExpr(p *parser) (astExpr, error) {
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

func parseInner(p *parser) (astExpr, error) {
	// last sym read was SYM_OPEN_PARENS and it was consumed
	// terminates with final ')' unconsumed
	sym := p.read()
	if isKnownOpType(sym) {
		return parseOp(p)
	}
	if sym.symType == SYM_SYM {
		return parseCallExpr(p)
	}
	if sym.symType == SYM_LET {
		return parseLetExpr(p)
	}
	return nil, fmt.Errorf("expected an operation or 'let', but got %s", sym.value)
}

func parseOp(p *parser) (astExpr, error) {
	// last sym read was SYM_OP and it is unconsumed
	// terminates with final ')' unconsumed
	sym := p.read()
	op := newOpTypeFromSym(sym)
	p.consume()
	args := []astExpr{}
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
	return newBinaryOpExpr(funcNameForBinOp(op), op, args[0], args[1]), nil
}

func parseLetExpr(p *parser) (astExpr, error) {
	// last sym read was SYM_LET and it is unconsumed
	// terminates successfully with final ')' still unconsumed
	p.consume() // eat 'let'
	sym := p.read()
	if sym.symType != SYM_OPEN_PAREN {
		return nil, fmt.Errorf("expected '(' at beginning of variable clause, but got %s", sym.value)
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

func parseFnStmt(p *parser) (*astFnStmt, error) {
	// last sym read was '(' and it is unconsumed
	// terminates successfully with final ')' still unconsumed
	sym := p.read()
	if sym.symType != SYM_OPEN_PAREN {
		return nil, fmt.Errorf("expected function def to start with '(' and not %q", sym.value)
	}
	p.consume() // eat '('

	sym = p.read()
	if sym.symType != SYM_FN {
		return nil, fmt.Errorf("expected function def to start with 'fn' and not %q", sym.value)
	}
	p.consume() // eat 'fn'

	// read new function name
	sym = p.read()
	if sym.symType != SYM_SYM {
		return nil, fmt.Errorf("expected a function name and not %q", sym.value)
	}
	p.consume() // eat function name
	name := sym.value

	// read parameter list
	sym = p.read()
	if sym.symType != SYM_OPEN_PAREN {
		return nil, fmt.Errorf("expected '(' at beginning of variable clause and not %q", sym.value)
	}
	p.consume() // eat '('


	// Read arguments
	p.symt.push()
	defer p.symt.pop()
	args := []string{}
argLoop:
	for true {
		sym = p.read()
		switch sym.symType {
		case SYM_SYM:
			p.consume()
			p.symt.add(sym.value, struct{}{})
			args = append(args, sym.value)
		case SYM_CLOSE_PAREN:
			p.consume()
			break argLoop
		default:
			return nil, fmt.Errorf("expecting variable or ')' and not %q", sym.value)
		}
	}

	// read function definition
	body, err := parseExpr(p)
	if err != nil {
		return nil, err
	}

	sym = p.read()
	if sym.symType != SYM_CLOSE_PAREN {
		return nil, fmt.Errorf("expected function def to end with ')' and not %q", sym.value)
	}
	p.consume() // eat ')'

	return newFnStmt(name, args, body), nil
}

func parseCallExpr(p *parser) (astExpr, error) {
	// last sym read was SYM_SYM and it is unconsumed
	// terminates successfully with final ')' still unconsumed

	// read function name
	sym := p.read()
	if sym.symType != SYM_SYM {
		return nil, fmt.Errorf("expected a function name and not %s", sym.value)
	}
	p.consume()
	name := sym.value

	// read arguments
	args := []astExpr{}
	for true {
		sym = p.read()
		if sym.symType == SYM_CLOSE_PAREN {
			break
		}
		arg, err := parseExpr(p)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	ce := newCallExpr(name, args)
	p.unresolved = append(p.unresolved, ce)
	return *ce, nil
}

func resolveFuncs(p *parser) error {
	for _, call := range p.unresolved {
		fn, ok := p.funcs[call.name]
		if !ok {
			return fmt.Errorf("function %q is not defined", call.name)
		}
		if len(call.args) != len(fn.args) {
			return fmt.Errorf(
				"function %s has %d args but was called with %d args",
				call.name,
				len(call.args),
				len(fn.args))
		}
	}
	return nil
}

type parser struct {
	symc       chan symbol
	readAhead  bool
	next       symbol
	astc       chan *astExpr
	symt       symTable
	unresolved []*astCallExpr
	funcs      map[string]*astFnStmt
}

func newParser(symc chan symbol) parser {
	return parser{
		symc:       symc,
		readAhead:  false,
		astc:       make(chan *astExpr),
		symt:       symTable{},
		unresolved: []*astCallExpr{},
		funcs:      map[string]*astFnStmt{},
	}
}

func (p *parser) read() symbol {
	if !p.readAhead {
		p.next = <-p.symc
		p.readAhead = true
	}
	return p.next
}

func (p *parser) consume() {
	if p.readAhead {
		p.readAhead = false
	}
}

func funcNameForBinOp(op opType) string {
	switch op {
	case OP_MINUS:
		return "minus"
	case OP_PLUS:
		return "plus"
	case OP_MULT:
		return "mult"
	case OP_DIV:
		return "div"
	default:
		panic("must implement func translation for binary op")
	}
}
