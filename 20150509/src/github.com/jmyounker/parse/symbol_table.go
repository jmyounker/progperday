package main

type symTable []map[string]interface{}

func (st *symTable)push() {
	*st = append(*st, map[string]interface{}{})
}

func (st *symTable)pop() {
	*st = (*st)[0:len(*st)-1]
}

func (st symTable)add(n string, v interface{}) {
	st[len(st)-1][n] = v
}

func (st symTable)get(n string) (interface{}, bool) {
	if len(st) == 0 {
		return struct{}{}, false
	}
	for i := len(st)-1; i >= 0; i-- {
		v, ok := st[i][n]
		if ok {
			return v, true
		}
	}
	return struct{}{}, false
}
