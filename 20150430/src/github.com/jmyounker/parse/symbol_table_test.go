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

func TestSymTableEmptyFrame(t *testing.T) {
	s := symTable{}
	s.push()
	_, ok := s.get("x")
	if ok {
		t.Fatalf("should not return a variable")
	}
}

func TestSymTableRetrieval(t *testing.T) {
	s := symTable{}
	s.push()
	s.add("x", "foo")
	v, ok := s.get("x")
	if !ok {
		t.Fatalf("should return var 'x'")
	}
	if v.(string) != "foo" {
		t.Fatalf("want foo, but got %s", v.(string))
	}
}

func TestSymTableLayeredRetrieval(t *testing.T) {
	s := symTable{}
	s.push()
	s.add("x", "foo")
	s.push()
	s.add("x", "bar")
	v, ok := s.get("x")
	if !ok {
		t.Fatalf("should return var 'x'")
	}
	if v.(string) != "bar" {
		t.Fatalf("want bar, but got %s", v.(string))
	}
}

func TestSymTableChecksLowerLayers(t *testing.T) {
	s := symTable{}
	s.push()
	s.add("x", "foo")
	s.push()
	s.add("y", "bar")
	v, ok := s.get("x")
	if !ok {
		t.Fatalf("should return var 'x'")
	}
	if v.(string) != "foo" {
		t.Fatalf("want foo, but got %s", v.(string))
	}
}
