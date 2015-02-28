package main

import (
	"sort"
	"strings"
	"testing"
	"reflect"
)

type E struct {n1, n2 int}

var inputGood = []struct{
	in string
	want []E
}{
	{"", []E{}},
	{"1 2", []E{E{1, 2}}},
	{"1 2\n2 3", []E{E{1, 2}, E{2, 3}}},
}

var inputBad = []struct{
	in string
	wantErrSubstr string
} {
	{"adsfd", "start"},
	{"43", "not enough"},
	{"43 asdf", "end"},
	{"43 42 41", "too many"},
}

func newGraph(edges []E) mapGraph {
	g := mapGraph{}
	for _, e := range edges {
		g.addEdge(e.n1, e.n2)
	}
	return g
}

func TestReadInputSuccesfully(t *testing.T) {
	for _, tc := range inputGood {
		g, err := readDataFrom(strings.NewReader(tc.in))
		assertNoError(t, err)
		if !reflect.DeepEqual(newGraph(tc.want), g) {
			t.Fatalf("want %#v but got %#v", newGraph(tc.want), g)
		}
	}
}

func TestReadInputUnuccesfully(t *testing.T) {
	for _, tc := range inputBad {
		_, err := readDataFrom(strings.NewReader(tc.in))
		assertError(t, err)
		if !strings.Contains(err.Error(), tc.wantErrSubstr) {
			t.Fatalf("expected no error to contain %q but was: %q", tc.wantErrSubstr, err)
		}
	}
}

func TestNodes(t *testing.T) {
	g := mapGraph{}
	assertEqualsUnordered(t, []int{}, g.nodes())
	g.addEdge(1, 2)
	assertEqualsUnordered(t, []int{1, 2}, g.nodes())
	g.addEdge(1, 3)
	assertEqualsUnordered(t, []int{1, 2, 3}, g.nodes())
}

func TestEdges(t *testing.T) {
	g := mapGraph{}
	assertEqualsUnordered(t, []int{}, g.nodes())
	g.addEdge(1, 2)
	assertEqualsUnordered(t, []int{2}, g.edges(1))
	assertEqualsUnordered(t, []int{1}, g.edges(2))
	g.addEdge(1, 3)
	assertEqualsUnordered(t, []int{2, 3}, g.edges(1))
	assertEqualsUnordered(t, []int{1}, g.edges(2))
	assertEqualsUnordered(t, []int{1}, g.edges(3))
}

func TestParallelEdges(t *testing.T) {
	g := mapGraph{}
	g.addEdge(1, 2)
	g.addEdge(1, 2)
	assertEqualsUnordered(t, []int{2, 2}, g.edges(1))
	assertEqualsUnordered(t, []int{1, 1}, g.edges(2))
}

func TestRemoveOneParallelEdge(t *testing.T) {
	g := mapGraph{}
	g.addEdge(1, 2)
	g.addEdge(1, 2)
	g.removeEdge(1, 2)
	assertEqualsUnordered(t, []int{2}, g.edges(1))
	assertEqualsUnordered(t, []int{1}, g.edges(2))
}

func TestRemoveEge(t *testing.T) {
	g := newGraph([]E{E{1, 2}, E{1, 3}})
	g.removeEdge(1, 3)
	assertEqualsUnordered(t, []int{1, 2, 3}, g.nodes())
	assertEqualsUnordered(t, []int{2}, g.edges(1))
	assertEqualsUnordered(t, []int{1}, g.edges(2))
	assertEqualsUnordered(t, []int{}, g.edges(3))
}

func TestRemoveNode(t *testing.T) {
	g := newGraph([]E{E{1, 2}, E{1, 3}})
	g.removeEdge(1, 3)
	g.removeNode(3)
	assertEqualGraph(t, newGraph([]E{E{1, 2}}), g)
}

func TestContractZero(t *testing.T) {
	g := mapGraph{}
	if (minConnPart(g) != 0) {
		t.Fatalf("expected no connections")
	}
}

func TestContractOne(t *testing.T) {
	g := newGraph([]E{E{1, 2}})
	g.removeEdge(1, 2)
	g.removeNode(2)
	if (minConnPart(g) != 0) {
		t.Fatalf("expected no connections")
	}
}

func TestContractTwo(t *testing.T) {
	g := newGraph([]E{E{1, 2}, E{1, 2}})
	if (minConnPart(g) != 2) {
		t.Fatalf("expected two connections")
	}
}

func TestContractOnceWithOneConnection(t *testing.T) {
	g := newGraph([]E{E{1, 2}, E{2, 3}})
	if (minConnPart(g) != 1) {
		t.Fatalf("expected one connections")
	}
}

func TestContractOnceWithTwoConnections(t *testing.T) {
	g := newGraph([]E{E{1, 2}, E{1, 2}, E{2, 3}, E{2, 3}})
	if (minConnPart(g) != 2) {
		t.Fatalf("expected one connections")
	}
}

func TestContractMultipleTimesWithTwoConnections(t *testing.T) {
	g := newGraph([]E{E{1, 2}, E{2, 3}, E{3, 1}})
	if (minConnPart(g) != 2) {
		t.Fatalf("expected two connections")
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected no error but got: %s", err)
	}
}

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Fatalf("Expected error but got: %s", err)
	}
}

func assertEquals(t *testing.T, want, got []int) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %#v but got %#v", want, got)
	}
}

func assertEqualGraph(t *testing.T, want, got mapGraph) {
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %#v but got %#v", want, got)
	}
}

func assertEqualsUnordered(t *testing.T, want, got []int) {
	sort.Sort(sort.IntSlice(want))
	sort.Sort(sort.IntSlice(got))
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %#v but got %#v", want, got)
	}
}
