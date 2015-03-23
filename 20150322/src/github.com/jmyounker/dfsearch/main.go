package main

import (
	"fmt"
	"strings"
)

func main() {
	g := mapDGraph{}

	g.addEdge(0, 1)
	g.addEdge(0, 4)
	g.addEdge(0, 5)
	g.addEdge(1, 4)
	g.addEdge(1, 7)
	g.addEdge(1, 2)
	g.addEdge(2, 3)
	g.addEdge(3, 0)
	g.addEdge(4, 6)

	p := breadthFirst(g, 0)

	s := []string{}
	for _, i := range p {
		s = append(s, fmt.Sprintf("%d", i))
	}
	fmt.Println(strings.Join(s, ", "))
}


func breadthFirst(g dGraph, s int) []int {
	path := []int{}
	visited := map[int]struct{}{}
	q := &Stack{}
	q.push(s)
	visited[s] = struct{}{}
	for !q.isEmpty() {
		n := q.pop()
		// visit the node
		path = append(path, n)

		// process the children
		for _, e := range g.edges(n) {
			_, ok := visited[e.end]
			if !ok {
				visited[e.end] = struct{}{}
				q.push(e.end)
			}
		}
	}
	return path
}


// Queue to manage breadth first search

type Stack []int

func (q *Stack) push(x int) {
	*q = append(*q, x)
}

func (q *Stack) pop() int {
	x := (*q)[len(*q)-1]
	*q = (*q)[0:len(*q)-1]
	return x
}

func (q *Stack) isEmpty() bool {
	return len(*q) == 0
}

