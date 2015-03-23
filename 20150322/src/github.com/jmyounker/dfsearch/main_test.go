package main

import (
	"testing"
)

func TestMinimalTree(t *testing.T) {
	g := mapDGraph{}
	g.addEdge(0, 1)

	assertEquals(t, []int{0, 1}, breadthFirst(g, 0))
}

func TestCircular(t *testing.T) {
	g := mapDGraph{}
	g.addEdge(0, 1)
	g.addEdge(1, 2)
	g.addEdge(2, 0)

	assertEquals(t, []int{0, 1, 2}, breadthFirst(g, 0))
}

func TestPairedTree(t *testing.T) {
	g := mapDGraph{}
	g.addEdge(0, 1)
	g.addEdge(0, 2)

	assertEquals(t, []int{0, 2, 1}, breadthFirst(g, 0))
}

func TestDepth(t *testing.T) {
	g := mapDGraph{}
	g.addEdge(0, 1)
	g.addEdge(1, 2)
	g.addEdge(0, 3)

	assertEquals(t, []int{0, 3, 1, 2}, breadthFirst(g, 0))
}

func TestDepthWithCircularLinks(t *testing.T) {
	g := mapDGraph{}
	g.addEdge(0, 1)
	g.addEdge(1, 2)
	g.addEdge(0, 3)
	g.addEdge(2, 0)
	g.addEdge(3, 1)

	assertEquals(t, []int{0, 3, 1, 2}, breadthFirst(g, 0))
}


func TestStackIsEmpty(t *testing.T) {
	s := Stack{}
	assertEquals(t, true, s.isEmpty())
}

func TestStackPushAndPull(t *testing.T) {
	s := Stack{}
	s.push(1)
	assertEquals(t, 1, s.pop())
	assertEquals(t, true, s.isEmpty())
}
