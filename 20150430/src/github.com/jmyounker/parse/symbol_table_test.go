package main

import (
	"testing"
)

func TestSymTablePushAndPop(t *testing.T) {
	s := symTable{}
	s.push()
	if len(s) != 1 {
		t.Fatalf("want len 1 and got len %d", len(s))
	}
	s.pop()
	if len(s) != 0 {
		t.Fatalf("want len 0 and got len %d", len(s))
	}
}

func TestSymTableGetEmpty(t *testing.T) {
	s := symTable{}
	_, ok := s.get("x")
	if ok {
		t.Fatalf("should not return a variable")
	}
}

