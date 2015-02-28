package main

import (
	"testing"
	"reflect"
)

var lexerTests = []struct{
	in string
	want []symbol
}{
	{"", []symbol{}},
	{"(", []symbol{symbol{SYM_OPEN_PAREN, "("}}},
	{")", []symbol{symbol{SYM_CLOSE_PAREN, ")"}}},
	{"*", []symbol{symbol{SYM_OP, "*"}}},
	{"/", []symbol{symbol{SYM_OP, "/"}}},
	{"+", []symbol{symbol{SYM_OP, "+"}}},
	{"-", []symbol{symbol{SYM_OP, "-"}}},
	{"-4", []symbol{symbol{SYM_LIT, "-4"}}},
	{"4", []symbol{symbol{SYM_LIT, "4"}}},
	{"42", []symbol{symbol{SYM_LIT, "42"}}},
	{"42(", []symbol{symbol{SYM_LIT, "42"}, symbol{SYM_OPEN_PAREN, "("}}},
	{"42)", []symbol{symbol{SYM_LIT, "42"}, symbol{SYM_CLOSE_PAREN, ")"}}},
	{"42 ", []symbol{symbol{SYM_LIT, "42"}}},
	{"x ", []symbol{symbol{SYM_ERR, "illegal symbol: x"}}},
	{"42x ", []symbol{symbol{SYM_ERR, "illegal symbol: x"}}},
	{" ( +  ) ", []symbol{
		symbol{SYM_OPEN_PAREN, "("},
		symbol{SYM_OP, "+"},
		symbol{SYM_CLOSE_PAREN, ")"},
	}},
}

func TestLexer(t *testing.T) {
	for _, tc := range lexerTests {
		symc := newLexer(tc.in)
		s := []symbol{}
		for sym := range symc {
			s = append(s, sym)
		}
		assertEquals(t, s, tc.want)
	}
}

func assertEquals(t *testing.T, want, got []symbol) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %#v but got %#v", want, got)
	}
}

