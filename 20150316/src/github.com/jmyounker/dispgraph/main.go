package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	g := mapDGraph{}

	g.addEdge(0, 1)
	g.addEdge(1, 2)
	g.addEdge(2, 3)
	g.addEdge(3, 0)

	g.addEdge(4, 5)
	g.addEdge(5, 6)

	parts := findPartitions(g);

	for _, p := range *parts {
		fmt.Println(p.String())
	}
}


type partition map[int]struct{}

type partID int
type partitionSet struct {
	// active partition
	part *partBuilder
	// the set of partitions discovered so far
	parts map[partID]*partBuilder
	// the mapping of nodes to partitions
	node2part map[int]*partBuilder
	// the set of unvisited nodes
	unvisited partition
	// the queue of nodes to visit next
	queue Queue
	// the current index of the partition being constructed
	partID partID
}

// Carries ID so we can easily locate which partition to remove
type partBuilder struct {
	id partID
	part *partition
}

func newPartBuilder(id partID) *partBuilder {
	return &partBuilder{id, &partition{}}
}

func (pb *partBuilder)add(node int) {
	(*pb.part)[node] = struct{}{}
}

func (pb *partBuilder)contains(node int) bool {
	_, ok := (*pb.part)[node]
	return ok
}

func (pb *partBuilder)String() string {
	return fmt.Sprintf("%d: %s", int(pb.id), pb.part.String())
}

func newPartitionSet() *partitionSet {
	return &partitionSet{
		part: nil,
		parts: map[partID]*partBuilder{},
		node2part: map[int]*partBuilder{},
		unvisited: partition{},
		queue: Queue{},
		partID: partID(0),
	}
}

func (ps *partitionSet)addNewPart() {
	ps.partID++
	ps.part = newPartBuilder(ps.partID)
	ps.parts[ps.partID] = ps.part
}

func (ps *partitionSet)deleteActivePart() {
	if ps.part == nil {
		return
	}
	delete(ps.parts, ps.part.id)
	ps.part = nil
}

func (ps *partitionSet)dump(){
	for i, p := range ps.parts {
		fmt.Printf("%d: %s\n", i, p.part.String())
	}
	if ps.part == nil {
		fmt.Println("act: nil")
	} else {
		fmt.Printf("act: %i: %s\n", ps.part.id, ps.part.part.String())
	}
}

func (ps partitionSet)pickUnvisisted() int {
	for ni := range ps.unvisited {
		return ni
	}
	panic("should never reach here")
}

func findPartitions(g dGraph) *[]*partition {
	ps := newPartitionSet()
	for _, n := range g.nodes() {
		ps.unvisited[n] = struct{}{}
	}

	// While there are unprocessed nodes
	for len(ps.unvisited) > 0 {
		ps.addNewPart()
		// Take first unvisited node. The loop condition guarantees one's existence.
		n := ps.pickUnvisisted();
		ps.queue.push(n)

		for !ps.queue.isEmpty() {
			n := ps.queue.pop()
			if processNode(n, g, ps) {
				for _, cn := range g.edges(n) {
					ps.queue.push(cn.end)
				}
			}
		}
	}
	p := []*partition{}
	for _, pb := range ps.parts {
		p = append(p, pb.part)
	}
	return &p
}

// True if this node has children
func processNode(node int, g dGraph, ps *partitionSet) bool {
	_, unvis := ps.unvisited[node]
	if unvis {
		delete(ps.unvisited, node)
		ps.part.add(node)
		ps.node2part[node] = ps.part
		return true
	}
	// This node has been seen in another partition, so this these are the same
	// partitions.
	inCurrent := ps.part.contains(node)
	other, inOther := ps.node2part[node]
	if inOther && !inCurrent {
		for pn := range *ps.part.part {
			other.add(pn)
			ps.node2part[pn] = other
		}
		ps.deleteActivePart()
		ps.part = other
	}
	return false
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
