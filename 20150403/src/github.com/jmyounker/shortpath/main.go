package main

import (
	"bufio"
	"container/heap"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)


var (
	fn = flag.String("f", "", "filename")
	sn = flag.Int("s", 0, "start node")
	en = flag.Int("e", 0, "end node")
)

func main() {
	flag.Parse()
	if *fn == "" {
		panic("must supply a filename")
	}

	g := readData(*fn)

	if !g.hasNode(*sn) {
		fmt.Fprintf(os.Stderr, "node %d not defined\n", *sn)
		os.Exit(127)
	}

	if !g.hasNode(*en) {
		fmt.Fprintf(os.Stderr, "node %d not defined\n", *en)
		os.Exit(127)
	}

	path, cost := shortestPath(g, *sn, *en)

	ps := []string{}
	for m := range path {
		ps = append(ps, fmt.Sprintf("%d", m))
	}
	fmt.Printf("%s: %d\n", strings.Join(ps, "->"), cost)
}

type edge struct {
	start int
	end int
	cost int
}

// mapDGraph is a directed graph indexed by edge
type mapDGraph map[int][]edge


type dGraph interface {
	addEdge(n1, n2, w int)
	hasNode(n int) bool
	nodes()[]int
	edges(n int)[]edge
	len() int
}


func (g mapDGraph) addEdge(start, end, cost int) {
	e := edge{start, end, cost}

	_, ok := g[start]
	if !ok {
		g[start] = []edge{}
	}

	_, ok = g[end]
	if !ok {
		g[end] = []edge{}
	}

	g[start] = append(g[start], e)
}

func (g mapDGraph)nodes() []int{
	n := []int{}
	for i := range g {
		n = append(n, i)
	}
	return n
}

func (g mapDGraph)hasNode(n int) bool {
	_, ok := g[n]
	return ok
}

func (g mapDGraph)edges(n int) []edge {
	es, ok := g[n]
	if !ok {
		panic("node must exist")
	}
	return es
}

func (g mapDGraph)len() int {
	return len(g)
}

func readData(fn string) dGraph {
	f, err := os.Open(fn)
	if err != nil {
		log.Fatalf("could not open %s for reading", fn)
	}
	defer f.Close()
	g, err := readDataFrom(f)
	if err != nil {
		panic(fmt.Sprintf("could not read file: %s", err))
	}
	return g
}

func readDataFrom(f io.Reader) (dGraph, error) {
	g := mapDGraph{}
	sf := bufio.NewScanner(f)
	sf.Split(bufio.ScanLines)
	for sf.Scan() {
		sl := bufio.NewScanner(strings.NewReader(sf.Text()))
		sl.Split(bufio.ScanWords)
		// if it's an empty line then we move on to the next line
		if (!sl.Scan()) {
			continue
		}

		// We have at least one field.
		n1, err := strconv.Atoi(sl.Text())
		if err != nil {
			return nil, fmt.Errorf("could not read edge start: %s", err)
		}

		// We expected a second field.
		if (!sl.Scan()) {
			return nil, errors.New("malformed input, not enough fields")
		}

		n2, err := strconv.Atoi(sl.Text())
		if err != nil {
			return nil, fmt.Errorf("could not read edge end: %s", err)
		}

		// We expected a third field.
		if (!sl.Scan()) {
			return nil, errors.New("malformed input, not enough fields")
		}

		w, err := strconv.Atoi(sl.Text())
		if err != nil {
			return nil, fmt.Errorf("could not read edge end: %s", err)
		}

		if w < 0 {
			return nil, fmt.Errorf("costs must be greater than zero")
		}

		if (sl.Scan()) {
			return nil, errors.New("malformed input, too many fields")
		}

		g.addEdge(n1, n2, w)
	}
	return g, nil
}


func shortestPath(g dGraph, start, end int) ([]int, int) {
	// Index to all the path costs.
	pci := map[int]*pathCost{}

	pci[start] = &pathCost{start, 0, math.MaxInt32, 0}
	eh := &costHeap{pci[start]}
	heap.Init(eh)

	// Calculate the lowest cost to each node and that nodes predecessor.
	for eh.Len() != 0 {
		pn := heap.Pop(eh).(*pathCost)
		for _, e := range g.edges(pn.node) {
			c := pn.cost + e.cost

			en, ok := pci[e.end]
			if ok {
				if c < en.cost {
					en.cost = c
					en.prev = pn.node
					heap.Fix(eh, en.index)
				}
			} else {
				en = &pathCost{e.end, c, pn.node, 0}
				pci[e.end] = en
				heap.Push(eh, en)
			}
		}
	}

	// Record the path from the end node back to the start.
	pn := pci[end]
	// If the end path is missing then the algorithm never found a path.
	if pn == nil {
		return []int{2}, math.MaxInt32
	}
	cost := pn.cost
	path_rev := []int{pn.node}
	for pn.prev != math.MaxInt32 {
		pn = pci[pn.prev]
		path_rev = append(path_rev, pn.node)
	}

	// Reverse that so we have the path from the start to the end.
	path := []int{}
	for i := len(path_rev)-1; i >= 0; i-- {
		path = append(path, path_rev[i])
	}

	return path, cost
}

type pathCost struct {
	node int // node number
	cost int // cost
	prev int // previous node producing this cost
	index int // location in heap store
}

type costHeap []*pathCost

func (eh costHeap) Len() int { return len(eh) }

func (eh costHeap) Less(i, j int) bool {
	return eh[i].cost < eh[j].cost
}
func (eh costHeap) Swap(i, j int) {
	eh[i], eh[j] = eh[j], eh[i]
	eh[i].index = i
	eh[j].index = j
}

func (eh *costHeap) Push(x interface{}) {
	n := len(*eh)
	pc := x.(*pathCost)
	pc.index = n
	*eh = append(*eh, pc)
}

func (eh *costHeap) Pop() interface{} {
	old := *eh
	n := len(old)
	pe := old[n-1]
	pe.index = -1 // for safety
	*eh = old[0 : n-1]
	return pe
}
