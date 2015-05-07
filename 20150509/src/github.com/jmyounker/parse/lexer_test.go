package main
import (
	"testing"
	"reflect"
)


var validTokens = []struct{
	in string
	want []symType
}{
	{"4", []symType{SYM_LIT}},
	{"-4", []symType{SYM_LIT}},
	{"- 4", []symType{SYM_OP, SYM_LIT}},
	{"-\t1\n", []symType{SYM_OP, SYM_LIT}},
	{"()", []symType{SYM_OPEN_PAREN, SYM_CLOSE_PAREN}},
	{"(<", []symType{SYM_OPEN_PAREN, SYM_OP}},
	{"* / + - <", []symType{SYM_OP, SYM_OP, SYM_OP, SYM_OP, SYM_OP}},
	{"if let fn foo", []symType{SYM_IF, SYM_LET, SYM_FN, SYM_SYM}},
}

func TestValidLexings(t *testing.T) {
	for _, tc := range validTokens {
		symc := newLexer(tc.in)
		syms := []symType{}
		for sym := range symc {
			syms = append(syms, sym.symType)
		}
		if !reflect.DeepEqual(tc.want, syms) {
			t.Fatalf("want %v but got %v", tc.want, syms)
		}
	}
}
