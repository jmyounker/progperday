package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)


var fn = flag.String("f", "", "filename")

func main() {
	flag.Parse()
	if *fn == "" {
		panic("must supply a filename")
	}
	g := readData(*fn)
	mcp := minConnPart(g)
	fmt.Printf("%d\n", mcp)
}

// Graph must support parallel edges, so needs a count
type mapGraph map[int]map[int]int


type graph interface {
	addEdge(n1, n2 int)
	removeEdge(n1, n2 int)
	removeNode(n int)
	nodes()[]int
	edges(n int)[]int
	len() int
}


func (g mapGraph) addEdge(n1, n2 int) {
	en1, ok := g[n1]
	if !ok {
		en1 = map[int]int{}
		g[n1] = en1
	}

	en2, ok := g[n2]
	if !ok {
		en2 = map[int]int{}
		g[n2] = en2
	}

	en1n2, ok := en1[n2]
	if !ok {
		en1n2 = 0
	}
	en1n2++
	en1[n2] = en1n2

	en2n1, ok := en2[n1]
	if !ok {
		en2n1 = 0
	}
	en2n1++;
	en2[n1] = en2n1
}

func (g mapGraph)removeEdge(n1, n2 int) {
	en1, ok := g[n1]
	if !ok {
		panic("can only remove defined edges")
	}

	en2, ok := g[n2]
	if !ok {
		panic("can only remove defined edges")
	}

	en1n2, ok := en1[n2]
	if !ok {
		panic("can only remove defined edges")
	}
	if en1n2 == 0 {
		panic("should never encounter an edge with count == 0")
	}

	en2n1, ok := en2[n1]
	if !ok {
		panic("can only remove defined edges")
	}
	if en2n1 == 0 {
		panic("should never encounter an edge with count == 0")
	}

	// Graph is undirected so cardinality should be the same from both sides
	if en1n2 != en2n1 {
		panic("mismatched number of inbound and outbound connections")
	}

	// We are now sure that both have the same cardinality, so we can just
	// use one edge.
	en1n2--

	if en1n2 == 0 {
		delete(en1, n2)
		delete(en2, n1)
	} else {
		en1[n2] = en1n2
		en2[n1] = en1n2
	}
}

func (g mapGraph)removeNode(n1 int) {
	en1, ok := g[n1]
	if !ok {
		panic("can only remove nodes that exist")
	}

	if len(en1) != 0 {
		panic("cannot delete nodes that still have edges")
	}

	delete(g, n1)
}

func (g mapGraph)nodes() []int{
	n := []int{}
	for i := range g {
		n = append(n, i)
	}
	return n
}

func (g mapGraph)edges(n int) []int {
	en, ok := g[n]
	if !ok {
		panic("node must exist")
	}
	m := []int{}
	for i, c := range en {
		for j := 0; j < c; j++ {
			m = append(m, i)
		}
	}
	return m
}

func (g mapGraph)len() int {
	return len(g)
}

func readData(fn string) graph {
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

func readDataFrom(f io.Reader) (graph, error) {
	g := mapGraph{}
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

		if (sl.Scan()) {
			return nil, errors.New("malformed input, too many fields")
		}

		g.addEdge(n1, n2)
	}
	return g, nil
}


func minConnPart(g graph) int {
	if g.len() == 0 {
		return 0
	}
	if g.len() == 1 {
		return 0
	}
	for {
		if g.len() == 2 {
			break
		}
		contract(g)
	}

	// number of connections between the two remaining nodes
	nds := g.nodes()
	return len(g.edges(nds[0]))
}

func contract(g graph) {
	nds := g.nodes()
	i := nds[rand.Intn(len(nds))]
	ie := g.edges(i)
	j := ie[rand.Intn(len(ie))]

	// move edges from j to i
	for _, k := range g.edges(j) {
		g.removeEdge(j, k)
		// eliminate loops to self
		if i != k {
			g.addEdge(i, k)
		}
	}
	// remove j
	g.removeNode(j)
}


