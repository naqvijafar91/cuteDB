package main

import (
	"reflect"
	"testing"
)

func TestAddElement(t *testing.T) {
	blockService := initBlockService()
	n, err := NewLeafNode([]int64{1, 2, 4, 5}, blockService)
	if err != nil {
		t.Error(err)
	}
	n.addElement(3)
	if !reflect.DeepEqual(n.getElements(), []int64{1, 2, 3, 4, 5}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}

	n, err = NewLeafNode([]int64{0}, blockService)
	if err != nil {
		t.Error(err)
	}
	n.addElement(3)
	if !reflect.DeepEqual(n.getElements(), []int64{0, 3}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}

	n, err = NewLeafNode([]int64{5, 10, 15}, blockService)
	if err != nil {
		t.Error(err)
	}
	n.addElement(3)
	if !reflect.DeepEqual(n.getElements(), []int64{3, 5, 10, 15}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}
}

func TestIsLeaf(t *testing.T) {
	blockService := initBlockService()
	child1, err := NewLeafNode([]int64{10, 11}, blockService)
	if err != nil {
		t.Error(err)
	}
	child2, err := NewLeafNode([]int64{13, 14}, blockService)
	if err != nil {
		t.Error(err)
	}
	n, err := NewNodeWithChildren([]int64{1, 2, 4, 5}, []uint64{child1.blockID, child2.blockID}, blockService)
	if err != nil {
		t.Error(err)
	}
	if n.isLeaf() {
		t.Error("Should not return as leaf as it has children", n)
	}

	child1, err = NewLeafNode(nil, blockService)
	if err != nil {
		t.Error(err)
	}
	child2, err = NewLeafNode(nil, blockService)
	if err != nil {
		t.Error(err)
	}
	n, err = NewNodeWithChildren([]int64{1, 2, 4, 5}, nil, blockService)
	if err != nil {
		t.Error(err)
	}
	if !n.isLeaf() {
		t.Error("Should return as leaf as it has no children", n)
	}

	n, err = NewLeafNode([]int64{1, 2, 4, 5}, blockService)
	if err != nil {
		t.Error(err)
	}
	if !n.isLeaf() {
		t.Error("Should return as leaf as it has no children", n)
	}
}

func TestHasOverFlown(t *testing.T) {
	blockService := initBlockService()
	elements := make([]int64, blockService.getMaxLeafSize()+1)
	for i := 0; i < blockService.getMaxLeafSize()+1; i++ {
		elements[i] = int64(i + 1)
	}
	n, err := NewLeafNode(elements, blockService)
	if err != nil {
		t.Error(err)
	}
	if !n.hasOverFlown() {
		t.Error("Should return true as node has overflown", n)
	}

	n, err = NewLeafNode([]int64{1, 2, 3}, blockService)
	if err != nil {
		t.Error(err)
	}
	if n.hasOverFlown() {
		t.Error("Should return false as node has 3 elements", n)
	}

	child1, err := NewLeafNode([]int64{10, 11}, blockService)
	if err != nil {
		t.Error(err)
	}
	child2, err := NewLeafNode([]int64{13, 14}, blockService)
	if err != nil {
		t.Error(err)
	}
	n, err = NewNodeWithChildren(elements, []uint64{child1.blockID,
		child2.blockID}, blockService)
	if err != nil {
		t.Error(err)
	}
	if !n.hasOverFlown() {
		t.Error("Should return true as node has overflown", n)
	}

}

func TestSplitLeafNode(t *testing.T) {
	blockService := initBlockService()
	n, err := NewLeafNode([]int64{1, 2, 3, 4, 5, 6}, blockService)
	if err != nil {
		t.Error(err)
	}
	poppedUpMiddleElement, leftChild, rightChild, err := n.SplitLeafNode()
	if err != nil {
		t.Error(err)
	}
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
	blockService := initBlockService()
	child1, err := NewLeafNode([]int64{10, 11, 12, 13, 14}, blockService)
	if err != nil {
		t.Error(err)
	}
	child2, err := NewLeafNode([]int64{30, 31, 32, 33, 34}, blockService)
	if err != nil {
		t.Error(err)
	}
	child3, err := NewLeafNode([]int64{40, 41, 42, 43, 44}, blockService)
	if err != nil {
		t.Error(err)
	}
	child4, err := NewLeafNode([]int64{50, 51, 52, 53, 54}, blockService)
	if err != nil {
		t.Error(err)
	}
	child5, err := NewLeafNode([]int64{60, 61, 62, 63, 64}, blockService)
	if err != nil {
		t.Error(err)
	}
	child6, err := NewLeafNode([]int64{70, 71, 72, 73, 74}, blockService)
	if err != nil {
		t.Error(err)
	}
	n, err := NewNodeWithChildren([]int64{15, 35, 45, 55, 65, 75}, []uint64{child1.blockID,
		child2.blockID, child3.blockID, child4.blockID, child5.blockID,
		child6.blockID}, blockService)
	if err != nil {
		t.Error(err)
	}
	poppedUpMiddleElement, leftChild, rightChild, err := n.SplitNonLeafNode()
	if err != nil {
		t.Error(err)
	}
	if poppedUpMiddleElement != 55 {
		t.Error("Wrong middle element, should be 55", poppedUpMiddleElement)
	}
	childToBeTested, err := leftChild.getChildAtIndex(2)
	if err != nil {
		t.Error(err)
	}
	if childToBeTested.getElementAtIndex(4) != 44 {
		t.Error("Element should be 44", childToBeTested.getElementAtIndex(4))
	}
	childToBeTested, err = leftChild.getChildAtIndex(3)
	if err != nil {
		t.Error(err)
	}
	if childToBeTested.getElementAtIndex(4) != 54 {
		t.Error("Element should be 54", childToBeTested.getElementAtIndex(4))
	}

	childToBeTested, err = rightChild.getChildAtIndex(1)
	if err != nil {
		t.Error(err)
	}
	if childToBeTested.getElementAtIndex(4) != 74 {
		t.Error("Element should be 44", childToBeTested.getElementAtIndex(4))
	}
}

func TestAddPoppedupElement(t *testing.T) {
	blockService := initBlockService()
	child1OfParent, err := NewLeafNode([]int64{1000, 1001, 1002, 1003, 1004}, blockService)
	if err != nil {
		t.Error(err)
	}
	child2OfParent, err := NewLeafNode([]int64{2000, 2001, 2002, 2003, 2004}, blockService)
	if err != nil {
		t.Error(err)
	}
	parentNode, err := NewNodeWithChildren([]int64{500}, []uint64{child1OfParent.blockID,
		child2OfParent.blockID}, blockService)
	child3, err := NewLeafNode([]int64{3000, 3001, 3002, 3003, 3004}, blockService)
	if err != nil {
		t.Error(err)
	}
	child4, err := NewLeafNode([]int64{4000, 4001, 4002, 4003, 4004}, blockService)
	if err != nil {
		t.Error(err)
	}
	parentNode.AddPoppedUpElementIntoCurrentNodeAndUpdateWithNewChildren(55, child3, child4)

	child, err := parentNode.getChildAtIndex(0)
	if err != nil {
		t.Error(err)
	}
	if child.getElementAtIndex(4) != 3004 {
		t.Error("Child not inserted at the correct position", child.getElements())
	}

}

