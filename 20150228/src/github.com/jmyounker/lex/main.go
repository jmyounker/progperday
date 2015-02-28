package main

// Lexer for a very simple list-like calculator language.  Next version
// will probably incorporate a parser, and then subsequently a compiler
// via llvm.

import (
	"flag"
	"fmt"
)

var prog = flag.String("p", "", "program")

func main() {
	flag.Parse()
	if *prog == "" {
		panic("must supply a program")
	}
	symc := newLexer(*prog)
	for sym := range symc {
		fmt.Printf("%#v\n", sym)
	}
}

type lexer struct {
	src string
	consumed int
	examining int
	symbol chan symbol
}

func newLexer(src string) chan symbol {
	sym := make(chan symbol)
	lxr := lexer{src, 0, 0, sym}
	go lxr.lex()
	return sym
}

func (lxr *lexer)lex() {
	state := startSymState
	for {
		state = state(lxr)
		if state == nil {
			close(lxr.symbol)
			return
		}
	}
}

func (lxr *lexer)current() string {
	return lxr.src[lxr.examining:lxr.examining+1]
}

func (lxr *lexer)advance() {
	lxr.examining++
}

func (lxr *lexer)consume() {
	lxr.consumed = lxr.examining
}

func (lxr *lexer)produce() string {
	res := lxr.src[lxr.consumed:lxr.examining]
	lxr.consume()
	return res
}

func (lxr *lexer)atEnd() bool {
	return lxr.examining >= len(lxr.src)
}

type symType int

const (
	SYM_OP symType = iota
	SYM_LIT
	SYM_OPEN_PAREN
	SYM_CLOSE_PAREN
	SYM_ERR
)

// A symbol is lexically valid string
type symbol struct {
	symType symType
	value string
}

func isDigit(c string) bool {
	return (c == "0" || c == "1" || c == "2" || c == "3" || c == "4" ||
		c == "5" || c == "6" || c == "7" || c == "8" || c == "9")
}

type stateFn func(*lexer) stateFn

func startSymState(lxr *lexer) stateFn {
	if lxr.atEnd() {
		return nil
	}
	c := lxr.current()
	if c == " " {
		lxr.advance()
		lxr.consume()
		return startSymState
	}
	if c == "(" {
		lxr.advance()
		lxr.symbol <- symbol{SYM_OPEN_PAREN, lxr.produce()}
		return startSymState
	}
	if c == ")" {
		lxr.advance()
		lxr.symbol <- symbol{SYM_CLOSE_PAREN, lxr.produce()}
		return startSymState
	}
	if c == "+" {
		lxr.advance()
		lxr.symbol <- symbol{SYM_OP, lxr.produce()}
		return startSymState
	}
	if c == "*"{
		lxr.advance()
		lxr.symbol <- symbol{SYM_OP, lxr.produce()}
		return startSymState
	}
	if c == "/" {
		lxr.advance()
		lxr.symbol <- symbol{SYM_OP, lxr.produce()}
		return startSymState
	}
	if c == "-" {
		lxr.advance()
		return negNumberState
	}
	if isDigit(c) {
		lxr.advance()
		return numberState
	}
	lxr.symbol <- symbol{SYM_ERR, fmt.Sprintf("illegal symbol: %s", c)}
	return nil
}

func negNumberState(lxr *lexer) stateFn {
	if lxr.atEnd() {
		lxr.symbol <- symbol{SYM_OP, lxr.produce()}
		return nil
	}
	c := lxr.current()
	if c == " " || c == "(" || c == ")" {
		lxr.symbol <- symbol{SYM_OP, lxr.produce()}
		return startSymState
	}
	if isDigit(c) {
		lxr.advance()
		return numberState
	}
	lxr.symbol <- symbol{SYM_ERR, fmt.Sprintf("illegal symbol: %s", c)}
	return nil
}

func numberState(lxr *lexer) stateFn {
	if lxr.atEnd() {
		lxr.symbol <- symbol{SYM_LIT, lxr.produce()}
		return nil
	}
	c := lxr.current()
	if c == " " || c == "(" || c == ")" {
		lxr.symbol <- symbol{SYM_LIT, lxr.produce()}
		return startSymState
	}
	if isDigit(c) {
		lxr.advance()
		return numberState
	}
	lxr.symbol <- symbol{SYM_ERR, fmt.Sprintf("illegal symbol: %s", c)}
	return nil
}
