package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	g := mapGraph{}

	g.addEdge(0, 1)
	g.addEdge(1, 2)
	g.addEdge(2, 3)

	g.addEdge(4, 5)
	g.addEdge(5, 6)

	parts := findPartitions(g);

	for _, p := range *parts {
		fmt.Println(p.String())
	}
}


type partition map[int]struct{}

func (p *partition)String() string {
	// Get list of keys and sort them.
	on := []int{}
	for n := range *p {
		on = append(on, n)
	}
	sort.Sort(sort.IntSlice(on))

	// Transform into comma separated string.
	s := []string{}
	for _, n := range on {
		s = append(s, fmt.Sprintf("%d", n))
	}
	return strings.Join(s, ", ")
}

type partitionSet struct {
	// active partition
	part *partition
	// the set of partitions discovered so far
	parts *[]*partition
	// the set of unvisited nodes
	unvisited partition
	// the queue of nodes to visit next
	queue Queue
}

func newPartitionSet() *partitionSet {
	return &partitionSet{
		part: nil,
		parts: &[]*partition{},
		unvisited: partition{},
		queue: Queue{},
	}
}

func (ps *partitionSet)addNewPart() {
	ps.part = &partition{}
	*ps.parts = append(*ps.parts, ps.part)
}

func (ps *partitionSet)dump(){
	for i, p := range *ps.parts {
		fmt.Printf("%d: %s\n", i, p.String())
	}
}

func (ps partitionSet)pickUnvisited() int {
	for ni := range ps.unvisited {
		return ni
	}
	panic("should never reach here")
}

func findPartitions(g graph) *[]*partition {
	ps := newPartitionSet()

	// Mark all as unvisited
	for _, n := range g.nodes() {
		ps.unvisited[n] = struct{}{}
	}

	// While there are unprocessed nodes
	for len(ps.unvisited) > 0 {
		ps.addNewPart()

		// Take first unvisited node. The loop condition guarantees one's existence.
		n := ps.pickUnvisited();
		ps.queue.push(n)

		for !ps.queue.isEmpty() {
			n := ps.queue.pop()
			_, unvis := ps.unvisited[n]
			if !unvis {
				continue
			}
			// Mark as visited
			delete(ps.unvisited, n)

			// Add to partition
			(*ps.part)[n] = struct{}{}

			// Examine children in subsequent passes
			for _, cn := range g.edges(n) {
				ps.queue.push(cn)
			}
		}
	}
	return ps.parts
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


