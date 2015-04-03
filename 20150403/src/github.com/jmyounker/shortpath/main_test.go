package main

import (
	"container/heap"
//	"fmt"
	"math"
	"reflect"
	"sort"
	"testing"
)

func TestGraphNodes(t *testing.T) {
	g := mapDGraph{}
	g.addEdge(1, 2, 1)

	assertEqualsUnordered(t, []int{1, 2}, g.nodes())
}

func TestGraphEdges(t *testing.T) {
	g := mapDGraph{}
	g.addEdge(1, 2, 1)
	g.addEdge(2, 1, 1)

	assertEquals(t, []edge{edge{1, 2, 1}}, g.edges(1))
	assertEquals(t, []edge{edge{2, 1, 1}}, g.edges(2))
}

func TestInitAndPop(t *testing.T) {
	p1 := &pathCost{1, 1, 0, 0}
	p2 := &pathCost{2, 2, 1, 0}
	p3 := &pathCost{3, 3, 2, 0}

	ch := &costHeap{p3, p1, p2}
	heap.Init(ch)

	e1 := heap.Pop(ch)
	e2 := heap.Pop(ch)
	e3 := heap.Pop(ch)

	assertEquals(t, costHeap{}, *ch)
	assertEquals(t, p1, e1)
	assertEquals(t, p2, e2)
	assertEquals(t, p3, e3)
}

func TestPushAndPop(t *testing.T) {
	ch := &costHeap{}
	heap.Init(ch)

	p1 := &pathCost{1, 1, 0, 0}
	p2 := &pathCost{2, 2, 0, 0}
	heap.Push(ch, p1)
	heap.Push(ch, p2)

	e1 := heap.Pop(ch)
	e2 := heap.Pop(ch)

	assertEquals(t, costHeap{}, *ch)
	assertEquals(t, p1, e1)
	assertEquals(t, p2, e2)
}


func TestInterfacePush(t *testing.T) {
	ch := &costHeap{}
	p1 := &pathCost{1, 1, 0, 0}
	p2 := &pathCost{2, 2, 0, 0}

	ch.Push(p1)
	ch.Push(p2)

	assertEquals(t, costHeap{p1, p2}, *ch)
}

func TestInterfacePop(t *testing.T) {
	p1 := &pathCost{1, 1, 0, 0}
	p2 := &pathCost{2, 2, 0, 0}
	ch := &costHeap{p1, p2}

	e2 := ch.Pop()
	e1 := ch.Pop()

	assertEquals(t, p1, e1)
	assertEquals(t, p2, e2)
	assertEquals(t, costHeap{}, *ch)
}

func TestPushAndPopOneItem(t *testing.T) {
	ch := &costHeap{}
	p1 := &pathCost{1, 3, 0, 0}

	heap.Push(ch, p1)
	e1 := heap.Pop(ch)

	assertEquals(t, p1, e1)
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

func TestGraphMinimalConnected(t *testing.T) {
	g := mapDGraph{}
	g.addEdge(1, 2, 1)

	path, cost := shortestPath(g, 1, 2)
	assertEquals(t, []int{1, 2}, path)
	assertEquals(t, 1, cost)
}

func TestGraphMinimalConnectedNoPath(t *testing.T) {
	g := mapDGraph{}
	g.addEdge(1, 2, 1)

	path, cost := shortestPath(g, 2, 1)
	assertEquals(t, []int{2}, path)
	assertEquals(t, math.MaxInt32, cost)
}

func TestGraphTraaversalStopsWhenExpectedNodeFound(t *testing.T) {
	g := mapDGraph{}
	g.addEdge(1, 2, 1)
	g.addEdge(2, 3, 1)

	path, cost := shortestPath(g, 1, 2)
	assertEquals(t, []int{1, 2}, path)
	assertEquals(t, 1, cost)
}

func TestGraphTraaversalOverMultipleLinks(t *testing.T) {
	g := mapDGraph{}
	g.addEdge(1, 2, 1)
	g.addEdge(2, 3, 1)

	path, cost := shortestPath(g, 1, 3)
	assertEquals(t, []int{1, 2, 3}, path)
	assertEquals(t, 2, cost)
}

func TestGraphTraaversalWithMultiplePaths(t *testing.T) {
	g := mapDGraph{}
	g.addEdge(1, 2, 1)
	g.addEdge(2, 3, 2)
	g.addEdge(1, 3, 4)
	g.addEdge(2, 4, 6)
	g.addEdge(3, 4, 3)

	path, cost := shortestPath(g, 1, 4)
	assertEquals(t, []int{1, 2, 3, 4}, path)
	assertEquals(t, 6, cost)
}
