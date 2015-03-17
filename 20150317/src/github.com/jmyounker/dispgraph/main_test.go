package main

import (
	"sort"
	"testing"
)


func TestPartitionToString(t *testing.T) {
	data := []struct{
		in []int
		want string
	}{
		{[]int{1}, "1"},
		{[]int{2, 1}, "1, 2"},
	}
	for _, tc := range data {
		p := partition{}
		for _, n := range tc.in {
			p[n] = struct{}{}
		}

		got := p.String()

		if got != tc.want {
			t.Fatal("want %q but got %q", tc.want, got)
		}
	}
}

func TestSmallestRecursiveGraph(t *testing.T) {
	g := mapGraph{}
	g.addEdge(1, 1)

	parts := findPartitions(g)

	if len(*parts) != 1 {
		t.Fatal("expected at least one")
	}
}

func TestNextSmallestRecursiveGraph(t *testing.T) {
	g := mapGraph{}
	g.addEdge(1, 2)
	g.addEdge(2, 1)

	parts := findPartitions(g)

	if len(*parts) != 1 {
		t.Fatalf("wanted only %d but got %d partitions", 1, len(*parts))
	}
	if (*parts)[0].String() != "1, 2" {
		t.Fatalf("wanted only partition %s but got %s", (*parts)[0].String(), "1, 2")
	}
}

func TestSmallestTreeShapedGraph(t *testing.T) {
	g := mapGraph{}
	g.addEdge(1, 2)
	g.addEdge(1, 3)

	parts := findPartitions(g)

	if len(*parts) != 1 {
		t.Fatalf("wanted only %d but got %d partitions", 1, len(*parts))
	}
	if (*parts)[0].String() != "1, 2, 3" {
		t.Fatalf("wanted only partition %s but got %s", (*parts)[0].String(), "1, 2, 3")
	}
}

func TestMultipleOutboundEdges(t *testing.T) {
	g := mapGraph{}
	g.addEdge(1, 3)
	g.addEdge(1, 3)

	parts := findPartitions(g)

	if len(*parts) != 1 {
		t.Fatalf("wanted only %d but got %d partitions", 1, len(*parts))
	}
	if (*parts)[0].String() != "1, 3" {
		t.Fatalf("wanted only partition %s but got %s", "1, 3", (*parts)[0].String())
	}
}

func TestTreeShapedGraph(t *testing.T) {
	g := mapGraph{}
	g.addEdge(1, 2)
	g.addEdge(1, 3)
	g.addEdge(3, 4)
	g.addEdge(3, 5)

	parts := findPartitions(g)

	if len(*parts) != 1 {
		t.Fatalf("wanted only %d but got %d partitions", 1, len(*parts))
	}
	if (*parts)[0].String() != "1, 2, 3, 4, 5" {
		t.Fatalf("wanted partition %s but got %s", "1, 2, 3, 4, 5", (*parts)[0].String())
	}
}

func TestSeparatePartitions(t *testing.T) {
	g := mapGraph{}
	// one partition
	g.addEdge(1, 2)
	g.addEdge(2, 3)
	// another partition
	g.addEdge(4, 5)
	g.addEdge(5, 6)

	parts := findPartitions(g)

	if len(*parts) != 2 {
		t.Fatalf("wanted only %d but got %d partitions", 2, len(*parts))
	}
	// make into a consistently ordered set
	got := []string{}
	for _, p := range *parts {
		got = append(got, (*p).String())
	}
	sort.Sort(sort.StringSlice(got))

	if got[0] != "1, 2, 3" {
		t.Fatalf("wanted partition %s but got %s", "1, 2, 3", got[0])
	}

	if got[1] != "4, 5, 6" {
		t.Fatalf("wanted partition %s but got %s", "4, 5, 6", got[1])
	}
}
