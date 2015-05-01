package main

// Lexer for a very simple list-like calculator language.  Next version
// will probably incorporate a parser, and then subsequently a compiler
// via llvm.

// TODO(jeff): Make silently consume CRLF

import (
	"fmt"
	"unicode/utf8"
)

type lexer struct {
	src string
	consumed int
	examining int
	rawSymbol chan symbol
	symbol chan symbol
}

func newLexer(src string) chan symbol {
	rawSym := make(chan symbol)
	sym := newKeywordTransformer(rawSym)
	lxr := lexer{src, 0, 0, rawSym, sym}
	go lxr.lex()
	return sym
}

func (lxr *lexer)lex() {
	state := startSymState
	for {
		state = state(lxr)
		if state == nil {
			close(lxr.rawSymbol)
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
	SYM_EOF symType = iota // defaults to EOF
	SYM_OP
	SYM_LIT
	SYM_OPEN_PAREN
	SYM_CLOSE_PAREN
	SYM_SYM
	SYM_LET
	SYM_FN
	SYM_ERR
)

// A symbol is a lexically valid string
type symbol struct {
	symType symType
	value string
}

func isDigit(c string) bool {
	r, _ := utf8.DecodeRuneInString(c[0:])
	return '0' <= r && r <= '9'
}

func isChar(c string) bool {
	r, _ := utf8.DecodeRuneInString(c[0:])
	if 'a' <= r && r <= 'z' {
		return true
	}
	if 'A' <= r && r <= 'Z' {
		return true
	}
	return false
}

type lxrStateFn func(*lexer) lxrStateFn

func startSymState(lxr *lexer) lxrStateFn {
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
		lxr.rawSymbol <- symbol{SYM_OPEN_PAREN, lxr.produce()}
		return startSymState
	}
	if c == ")" {
		lxr.advance()
		lxr.rawSymbol <- symbol{SYM_CLOSE_PAREN, lxr.produce()}
		return startSymState
	}
	if c == "+" {
		lxr.advance()
		lxr.rawSymbol <- symbol{SYM_OP, lxr.produce()}
		return startSymState
	}
	if c == "*"{
		lxr.advance()
		lxr.rawSymbol <- symbol{SYM_OP, lxr.produce()}
		return startSymState
	}
	if c == "/" {
		lxr.advance()
		lxr.rawSymbol <- symbol{SYM_OP, lxr.produce()}
		return startSymState
	}
	if c == "-" {
		lxr.advance()
		return negNumberState
	}
	if isChar(c) {
		lxr.advance()
		return symState
	}
	if isDigit(c) {
		lxr.advance()
		return numberState
	}
	lxr.rawSymbol <- symbol{SYM_ERR, fmt.Sprintf("illegal symbol: %s", c)}
	return nil
}

func symState(lxr *lexer) lxrStateFn {
	if lxr.atEnd() {
		lxr.rawSymbol <- symbol{SYM_SYM, lxr.produce()}
		return nil
	}
	c := lxr.current()
	if c == " " || c == "(" || c == ")" {
		lxr.rawSymbol <- symbol{SYM_SYM, lxr.produce()}
		return startSymState
	}
	if isChar(c) || isDigit(c) {
		lxr.advance()
		return symState
	}
	lxr.rawSymbol <- symbol{SYM_ERR, fmt.Sprintf("illegal symbol: %s", c)}
	return nil
}

func negNumberState(lxr *lexer) lxrStateFn {
	if lxr.atEnd() {
		lxr.rawSymbol <- symbol{SYM_OP, lxr.produce()}
		return nil
	}
	c := lxr.current()
	if c == " " || c == "(" || c == ")" {
		lxr.rawSymbol <- symbol{SYM_OP, lxr.produce()}
		return startSymState
	}
	if isDigit(c) {
		lxr.advance()
		return numberState
	}
	lxr.rawSymbol <- symbol{SYM_ERR, fmt.Sprintf("illegal symbol: %s", c)}
	return nil
}

func numberState(lxr *lexer) lxrStateFn {
	if lxr.atEnd() {
		lxr.rawSymbol <- symbol{SYM_LIT, lxr.produce()}
		return nil
	}
	c := lxr.current()
	if c == " " || c == "(" || c == ")" {
		lxr.rawSymbol <- symbol{SYM_LIT, lxr.produce()}
		return startSymState
	}
	if isDigit(c) {
		lxr.advance()
		return numberState
	}
	lxr.rawSymbol <- symbol{SYM_ERR, fmt.Sprintf("illegal symbol: %s", c)}
	return nil
}

// keywordLexer wraps a plain lexer and translates symType SYM_SYM into more specific
// symbols
type keywordTransformer struct {
	symIn chan symbol
	symOut chan symbol
	terminate chan struct{}
}

func newKeywordTransformer(symIn chan symbol) chan symbol {
	kwt := keywordTransformer{
		symIn: symIn,
		symOut: make(chan symbol),
	}
	go kwt.transformKeywords()
	return kwt.symOut
}

func (kwt *keywordTransformer)transformKeywords() {
	for sym := range kwt.symIn {
		if sym.symType != SYM_SYM {
			kwt.symOut <- sym
			continue
		}
		switch sym.value {
		case "let":
			sym.symType = SYM_LET
		case "fn":
			sym.symType = SYM_FN
		default:
			// do nothing
		}
		kwt.symOut <- sym
	}
	close(kwt.symOut)
}
