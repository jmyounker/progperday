package main

type edge struct {
	start int
	end int
	weight int
}

// mapDGraph is a directed graph indexed by node
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
