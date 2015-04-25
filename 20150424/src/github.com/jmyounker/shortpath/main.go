package main

import (
	"fmt"
	"math"
	"strings"
)


func main() {
	g := mapDGraph{}
	g.addEdge(1, 3, -2)
	g.addEdge(3, 4, 2)
	g.addEdge(4, 2, -1)
	g.addEdge(2, 1, 4)
	g.addEdge(2, 3, 3)

	m := calculatePathMap(g)
	path, cost := shortestPath(m, 4, 1)

	ps := []string{}
	for _, m:= range path {
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

func calculatePathMap(g dGraph) pathCostMap {
	m := pathCostMap{}

	for _, i := range (g.nodes()) {
		for _, j := range (g.nodes()) {
			m.put(i, j, &pathCost{math.MaxInt32, math.MaxInt32})
		}
	}

	for _, i := range (g.nodes()) {
		m.put(i, i, &pathCost{0, math.MaxInt32})
	}

	for _, i := range (g.nodes()) {
		for _, e := range (g.edges(i)) {
			m.put(e.start, e.end, &pathCost{e.cost, e.end})
		}
	}

	for _, k := range (g.nodes()) {
		for _, i := range (g.nodes()) {
			for _, j := range (g.nodes()) {
				pcc := m.get(i, j)
				pc1 := m.get(i, k)
				pc2 := m.get(k, j)
				// if no path then do nothing
				if pc1.cost == math.MaxInt32 || pc2.cost == math.MaxInt32 {
					continue
				}
				nc := pc1.cost + pc2.cost
				if nc >= pcc.cost {
					continue
				}
				pcc.cost = nc
				pcc.next = pc1.next
			}
		}
	}
	return m
}

func shortestPath(m pathCostMap, start, end int) ([]int, int){
	// Record the path from the end node back to the start.
	p := m.get(start, end)
	cost := p.cost
	path := []int{start}
	if p.cost == math.MaxInt32 {
		return path, cost
	}
	for p.next != end {
		path = append(path, p.next)
		p = m.get(p.next, end)
	}
	path = append(path, p.next)
	return path, cost
}


type pathCostMap map[int]map[int]*pathCost

func (m pathCostMap)get(x, y int)*pathCost {
	my, ok := m[x]
	if !ok {
		panic("should never access something that has not been written")
	}
	pc, ok := my[y]
	if !ok {
		panic("should never access something that has not been written")
	}
	return pc
}

func (m pathCostMap)put(x, y int, pc* pathCost) {
	my, ok := m[x]
	if !ok {
		my = map[int]*pathCost{}
		m[x] = my
	}
	my[y] = pc
}

type pathCost struct {
	cost int // cost
	next int // next node
}
