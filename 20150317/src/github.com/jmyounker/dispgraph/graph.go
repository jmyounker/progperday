package main

// Graph must support parallel edges, so needs a count
type mapGraph map[int]map[int]int


type graph interface {
	addEdge(n1, n2 int)
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

