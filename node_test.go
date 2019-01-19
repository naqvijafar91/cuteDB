package main

import (
	"reflect"
	"testing"
)

func TestAddElement(t *testing.T) {
	n := NewLeafNode([]int64{1, 2, 4, 5})
	n.addElement(3)
	if !reflect.DeepEqual(n.getElements(), []int64{1, 2, 3, 4, 5}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}

	n = NewLeafNode([]int64{0})
	n.addElement(3)
	if !reflect.DeepEqual(n.getElements(), []int64{0, 3}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}

	n = NewLeafNode([]int64{5, 10, 15})
	n.addElement(3)
	if !reflect.DeepEqual(n.getElements(), []int64{3, 5, 10, 15}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}
}

func TestIsLeaf(t *testing.T) {
	child1 := NewLeafNode([]int64{10, 11})
	child2 := NewLeafNode([]int64{13, 14})
	n := NewNodeWithChildren([]int64{1, 2, 4, 5}, []*Node{child1, child2})
	if n.isLeaf() {
		t.Error("Should not return as leaf as it has children", n)
	}

	child1 = NewLeafNode(nil)
	child2 = NewLeafNode(nil)
	n = NewNodeWithChildren([]int64{1, 2, 4, 5}, nil)
	if !n.isLeaf() {
		t.Error("Should return as leaf as it has no children", n)
	}

	n = NewLeafNode([]int64{1, 2, 4, 5})
	if !n.isLeaf() {
		t.Error("Should return as leaf as it has no children", n)
	}
}

func TestHasOverFlown(t *testing.T) {
	n := NewLeafNode([]int64{1, 2, 3, 4, 5, 6})
	if !n.hasOverFlown() {
		t.Error("Should return true as node has 6 elements", n)
	}

	n = NewLeafNode([]int64{1, 2, 3})
	if n.hasOverFlown() {
		t.Error("Should return false as node has 3 elements", n)
	}

	child1 := NewLeafNode([]int64{10, 11})
	child2 := NewLeafNode([]int64{13, 14})
	n = NewNodeWithChildren([]int64{1, 2, 3, 4, 5, 6}, []*Node{child1, child2})
	if !n.hasOverFlown() {
		t.Error("Should return true as node has 6 elements", n)
	}

}

func TestSplitLeafNode(t *testing.T) {
	n := NewLeafNode([]int64{1, 2, 3, 4, 5, 6})
	poppedUpMiddleElement, leftChild, rightChild := n.SplitLeafNode()
	if poppedUpMiddleElement != 4 {
		t.Error("Wrong middle Element popped up", poppedUpMiddleElement)
	}
	if leftChild.getElementAtIndex(2) != 3 {
		t.Error("Wrong value at leftchild", leftChild)
	}
	if rightChild.getElementAtIndex(1) != 6 {
		t.Error("Wrong value at rightchild ", rightChild)
	}
}

func TestSplitNonLeafNode(t *testing.T) {
	child1 := NewLeafNode([]int64{10, 11, 12, 13, 14})
	child2 := NewLeafNode([]int64{30, 31, 32, 33, 34})
	child3 := NewLeafNode([]int64{40, 41, 42, 43, 44})
	child4 := NewLeafNode([]int64{50, 51, 52, 53, 54})
	child5 := NewLeafNode([]int64{60, 61, 62, 63, 64})
	child6 := NewLeafNode([]int64{70, 71, 72, 73, 74})
	n := NewNodeWithChildren([]int64{15, 35, 45, 55, 65, 75}, []*Node{child1, child2, child3,
		child4, child5, child6})
	poppedUpMiddleElement, leftChild, rightChild := n.SplitNonLeafNode()
	if poppedUpMiddleElement != 55 {
		t.Error("Wrong middle element, should be 55", poppedUpMiddleElement)
	}
	if leftChild.getChildAtIndex(2).getElementAtIndex(4) != 44 {
		t.Error("Element should be 44", leftChild.getChildAtIndex(2).getElementAtIndex(4))
	}
	if leftChild.getChildAtIndex(3).getElementAtIndex(4) != 54 {
		t.Error("Element should be 54", leftChild.getChildAtIndex(3).getElementAtIndex(4))
	}
	if rightChild.getChildAtIndex(1).getElementAtIndex(4) != 74 {
		t.Error("Element should be 44", rightChild.getChildAtIndex(2).getElementAtIndex(4))
	}
}

func TestAddPoppedupElement(t *testing.T) {

	child1OfParent := NewLeafNode([]int64{1000, 1001, 1002, 1003, 1004})
	child2OfParent := NewLeafNode([]int64{2000, 2001, 2002, 2003, 2004})
	parentNode := NewNodeWithChildren([]int64{500}, []*Node{child1OfParent, child2OfParent})
	child3 := NewLeafNode([]int64{3000, 3001, 3002, 3003, 3004})
	child4 := NewLeafNode([]int64{4000, 4001, 4002, 4003, 4004})
	parentNode.AddPoppedUpElementIntoCurrentNodeAndUpdateWithNewChildren(55, child3, child4)
	if parentNode.getChildAtIndex(0).getElementAtIndex(4) != 3004 {
		t.Error("Child not inserted at the correct position", parentNode.getChildAtIndex(0).getElements())
	}

}
