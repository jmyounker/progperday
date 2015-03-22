package main

import (
	"fmt"
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

	breadthFirst(g, 0);
}


func breadthFirst(g dGraph, s int)  {
	visited := map[int]struct{}{}
	q := &Queue{}
	q.push(s)
	visited[s] = struct{}{}
	for !q.isEmpty() {
		n := q.pop()
		// visit the node
		fmt.Printf("%d\n", n)

		// process the children
		for _, e := range g.edges(n) {
			_, ok := visited[e.end]
			if !ok {
				visited[e.end] = struct{}{}
				q.push(e.end)
			}
		}
	}
}


// Queue to manage breadth first search

type Queue []int

func (q *Queue) push(x int) {
	*q = append(*q, x)
}

func (q *Queue) pop() int {
	x := (*q)[0]
	*q = (*q)[1:]
	return x
}

func (q *Queue) isEmpty() bool {
	return len(*q) == 0
}

