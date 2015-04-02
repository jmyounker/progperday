//package main
//
//import (
//	"reflect"
//	"testing"
//)
//
//func TestPairedTree(t *testing.T) {
//	g := mapDGraph{}
//	g.addEdge(0, 1)
//	g.addEdge(0, 2)
//
//	assertEquals(t, []int{0, 1, 2}, breadthFirst(g, 0))
//}
//
//func TestDepth(t *testing.T) {
//	g := mapDGraph{}
//	g.addEdge(0, 1)
//	g.addEdge(1, 2)
//	g.addEdge(0, 3)
//
//	assertEquals(t, []int{0, 1, 2, 3}, breadthFirst(g, 0))
//}
//
//func assertEquals(t *testing.T, want, got interface{}) {
//	if !reflect.DeepEqual(want, got) {
//		t.Fatalf("want %#v but got %#v", want, got)
//	}
//}
//
