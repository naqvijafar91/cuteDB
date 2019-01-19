package main

import (
	"reflect"
	"testing"
)

func TestAddElement(t *testing.T) {
	n := NewNodeWithoutChildren([]int64{1, 2, 4, 5})
	n.addElement(3)
	if !reflect.DeepEqual(n.getElements(), []int64{1, 2, 3, 4, 5}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}

	n = NewNodeWithoutChildren([]int64{0})
	n.addElement(3)
	if !reflect.DeepEqual(n.getElements(), []int64{0, 3}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}

	n = NewNodeWithoutChildren([]int64{5, 10, 15})
	n.addElement(3)
	if !reflect.DeepEqual(n.getElements(), []int64{3, 5, 10, 15}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}
}

func TestIsLeaf(t *testing.T) {
	child1 := NewNodeWithoutChildren([]int64{10, 11})
	child2 := NewNodeWithoutChildren([]int64{13, 14})
	n := NewNodeWithChildren([]int64{1, 2, 4, 5}, []*Node{child1, child2})
	if n.isLeaf() {
		t.Error("Should not return as leaf as it has children", n)
	}

	child1 = NewNodeWithoutChildren(nil)
	child2 = NewNodeWithoutChildren(nil)
	n = NewNodeWithChildren([]int64{1, 2, 4, 5}, nil)
	if !n.isLeaf() {
		t.Error("Should return as leaf as it has no children", n)
	}

	n = NewNodeWithoutChildren([]int64{1, 2, 4, 5})
	if !n.isLeaf() {
		t.Error("Should return as leaf as it has no children", n)
	}
}

func TestHasOverFlown(t *testing.T) {
	n := NewNodeWithoutChildren([]int64{1, 2, 3, 4, 5, 6})
	if !n.hasOverFlown() {
		t.Error("Should return true as node has 6 elements", n)
	}

	n = NewNodeWithoutChildren([]int64{1, 2, 3})
	if n.hasOverFlown() {
		t.Error("Should return false as node has 3 elements", n)
	}

}
