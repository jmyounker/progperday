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
	weight int
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


func (g mapDGraph) addEdge(start, end, weight int) {
	e := edge{start, end, weight}

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
			return nil, fmt.Errorf("weights must be greater than zero")
		}

		if (sl.Scan()) {
			return nil, errors.New("malformed input, too many fields")
		}

		g.addEdge(n1, n2, w)
	}
	return g, nil
}


func shortestPath(g dGraph, start, end int) ([]int, int) {
	// Allows us to locate all remaining nodes in the heap
	rem := map[int]*pathCost{}

	// Edges with lowest cost
	eh := &costHeap{}

	for _, n := range g.nodes() {
		pe := &pathCost{n, math.MaxInt32, 0}
		rem[n] = pe
		*eh = append(*eh, pe)
	}

	path := []int{}
	cost := 0
	rem[start].cost = cost
	heap.Init(eh)

	for {
		cn := heap.Pop(eh).(*pathCost)
		cost = cn.cost
		// no remaining path found
		if cn.cost == math.MaxInt32 {
			break
		}

		path = append(path, cn.node)
		delete(rem, cn.node)

		// reached end node
		if cn.node == end {
			break
		}

		// update costs to each edge
		for _, e := range g.edges(cn.node) {
			pc, ok := rem[e.end]
			if ok {
//				fmt.Printf("NODE: %#v; %d\n", cn, pc.cost)
				pc.cost = cn.cost + e.weight
				heap.Fix(eh, pc.index)
			}
		}
	}

	return path, cost
}

type pathCost struct {
	node int
	cost int
	index int
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
