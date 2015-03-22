package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestGraphNodes(t *testing.T) {
	g := mapDGraph{}
	g.addEdge(1, 2)

	assertEqualsUnordered(t, []int{1, 2}, g.nodes())
}

func TestGraphEdges(t *testing.T) {
	g := mapDGraph{}
	g.addEdge(1, 2)
	g.addEdge(2, 1)

	assertEquals(t, []edge{edge{1, 2}}, g.edges(1))
	assertEquals(t, []edge{edge{2, 1}}, g.edges(2))
}

func assertEquals(t *testing.T, want, got interface{}) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %#v but got %#v", want, got)
	}
}

func assertEqualsUnordered(t *testing.T, want, got []int) {
	sort.Sort(sort.IntSlice(want))
	sort.Sort(sort.IntSlice(got))
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %#v but got %#v", want, got)
	}
}
