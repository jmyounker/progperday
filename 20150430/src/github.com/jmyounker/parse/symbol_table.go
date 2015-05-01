package main

import "github.com/axw/gollvm/llvm"

type symTable []map[string]llvm.Value

func (st *symTable)push() {
	*st = append(*st, map[string]llvm.Value{})
}

func (st *symTable)pop() {
	*st = (*st)[0:len(*st)-1]
}

func (st symTable)add(n string, v llvm.Value) {
	st[len(st)-1][n] = v
}

func (st symTable)get(n string) (llvm.Value, bool) {
	if len(st) == 0 {
		return llvm.Value{}, false
	}
	for i := len(st)-1; i >= 0; i-- {
		v, ok := st[i][n]
		if ok {
			return v, true
		}
	}
	return llvm.Value{}, false
}
