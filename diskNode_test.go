package main

import (
	"fmt"
	"reflect"
	"testing"
)

func printNodeElements(n *DiskNode) {
	for i := 0; i < len(n.getElements()); i++ {
		fmt.Println(n.getElementAtIndex(i).key, n.getElementAtIndex(i).value)
	}
}
func TestAddElement(t *testing.T) {
	blockService := initBlockService()
	elements := make([]*Pairs, 3)
	elements[0] = NewPair("hola", "amigos")
	elements[1] = NewPair("foo", "bar")
	elements[2] = NewPair("gooz", "bumps")
	n, err := NewLeafNode(elements, blockService)
	if err != nil {
		t.Error(err)
	}
	addedElement := NewPair("added", "please check")
	n.addElement(addedElement)

	if !reflect.DeepEqual(n.getElements(), []*Pairs{addedElement, elements[0],
		elements[1], elements[2]}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}

	n, err = NewLeafNode([]*Pairs{NewPair("first", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	n.addElement(NewPair("second", "value"))
	if !reflect.DeepEqual(n.getElements(), []*Pairs{NewPair("first", "value"),
		NewPair("second", "value")}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}

	n, err = NewLeafNode([]*Pairs{NewPair("first", "value"),
		NewPair("second", "value"), NewPair("third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	n.addElement(NewPair("fourth", "value"))
	if !reflect.DeepEqual(n.getElements(), []*Pairs{NewPair("first", "value"),
		NewPair("fourth", "value"), NewPair("second", "value"), NewPair("third", "value")}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}
}

func TestIsLeaf(t *testing.T) {
	blockService := initBlockService()
	child1, err := NewLeafNode([]*Pairs{NewPair("first", "value"),
		NewPair("second", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child2, err := NewLeafNode([]*Pairs{NewPair("third", "value"),
		NewPair("forth", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	n, err := NewNodeWithChildren([]*Pairs{NewPair("fifth", "value"),
		NewPair("sixth", "value")}, []uint64{child1.blockID, child2.blockID}, blockService)
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
	n, err = NewNodeWithChildren([]*Pairs{NewPair("first", "value"),
		NewPair("second", "value")}, nil, blockService)
	if err != nil {
		t.Error(err)
	}
	if !n.isLeaf() {
		t.Error("Should return as leaf as it has no children", n)
	}

	n, err = NewLeafNode([]*Pairs{NewPair("first", "value"),
		NewPair("second", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	if !n.isLeaf() {
		t.Error("Should return as leaf as it has no children", n)
	}
}

func TestHasOverFlown(t *testing.T) {
	blockService := initBlockService()
	elements := make([]*Pairs, blockService.getMaxLeafSize()+1)
	for i := 0; i < blockService.getMaxLeafSize()+1; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		elements[i] = NewPair(key, value)
	}
	n, err := NewLeafNode(elements, blockService)
	if err != nil {
		t.Error(err)
	}
	if !n.hasOverFlown() {
		t.Error("Should return true as node has overflown", n)
	}

	n, err = NewLeafNode([]*Pairs{NewPair("first", "value"), NewPair("fourth", "value"),
		NewPair("second", "value"), NewPair("third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	if n.hasOverFlown() {
		t.Error("Should return false as node has 3 elements", n)
	}

	child1, err := NewLeafNode([]*Pairs{NewPair("first", "value"),
		NewPair("second", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child2, err := NewLeafNode([]*Pairs{NewPair("third", "value"),
		NewPair("fourth", "value")}, blockService)
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
	n, err := NewLeafNode([]*Pairs{NewPair("first", "value"),
		NewPair("fourth", "value"), NewPair("second", "value"), NewPair("third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	poppedUpMiddleElement, leftChild, rightChild, err := n.SplitLeafNode()
	if err != nil {
		t.Error(err)
	}
	if poppedUpMiddleElement.key != "second" {
		t.Error("Wrong middle Element popped up", poppedUpMiddleElement)
	}
	if leftChild.getElementAtIndex(1).key != "fourth" {
		t.Error("Wrong value at leftchild", leftChild)
	}
	if rightChild.getElementAtIndex(0).key != "third" {
		t.Error("Wrong value at rightchild ", rightChild)
	}
}

func TestSplitNonLeafNode(t *testing.T) {
	blockService := initBlockService()
	child1, err := NewLeafNode([]*Pairs{NewPair("1first", "value"),
		NewPair("1fourth", "value"), NewPair("1second", "value"), NewPair("1third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child2, err := NewLeafNode([]*Pairs{NewPair("2first", "value"),
		NewPair("2fourth", "value"), NewPair("2second", "value"), NewPair("2third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child3, err := NewLeafNode([]*Pairs{NewPair("3first", "value"),
		NewPair("3fourth", "value"), NewPair("3second", "value"), NewPair("3third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child4, err := NewLeafNode([]*Pairs{NewPair("4first", "value"),
		NewPair("4fourth", "value"), NewPair("4second", "value"), NewPair("4third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child5, err := NewLeafNode([]*Pairs{NewPair("5first", "value"),
		NewPair("5fourth", "value"), NewPair("5second", "value"), NewPair("5third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}

	n, err := NewNodeWithChildren([]*Pairs{NewPair("nfirst", "value"),
		NewPair("nfourth", "value"), NewPair("nsecond", "value"), NewPair("nthird", "value")},
		[]uint64{child1.blockID, child2.blockID, child3.blockID,
			child4.blockID, child5.blockID}, blockService)
	if err != nil {
		t.Error(err)
	}
	poppedUpMiddleElement, leftChild, rightChild, err := n.SplitNonLeafNode()
	if err != nil {
		t.Error(err)
	}
	if poppedUpMiddleElement.key != "nsecond" {
		t.Error("Wrong middle element, should be second", poppedUpMiddleElement)
	}
	childToBeTested, err := leftChild.getChildAtIndex(2)
	if err != nil {
		t.Error(err)
	}
	if childToBeTested.getElementAtIndex(2).key != "3second" {
		t.Error("Element should be 3second", childToBeTested.getElementAtIndex(2).key)
	}
	childToBeTested, err = leftChild.getChildAtIndex(1)
	if err != nil {
		t.Error(err)
	}
	if childToBeTested.getElementAtIndex(3).key != "2third" {
		t.Error("Element should be 2third", childToBeTested.getElementAtIndex(3).key)
	}

	childToBeTested, err = rightChild.getChildAtIndex(1)
	if err != nil {
		t.Error(err)
	}
	if childToBeTested.getElementAtIndex(3).key != "5third" {
		t.Error("Element should be 5third", childToBeTested.getElementAtIndex(3).key)
	}
}

func TestAddPoppedupElement(t *testing.T) {
	blockService := initBlockService()
	child1OfParent, err := NewLeafNode([]*Pairs{NewPair("1first", "value"),
	NewPair("1fourth", "value"), NewPair("1second", "value"), NewPair("1third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child2OfParent, err := NewLeafNode([]*Pairs{NewPair("2first", "value"),
	NewPair("2fourth", "value"), NewPair("2second", "value"), NewPair("2third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	parentNode, err := NewNodeWithChildren([]*Pairs{NewPair("parentfirst", "value")}, []uint64{child1OfParent.blockID,
		child2OfParent.blockID}, blockService)
	child3, err := NewLeafNode([]*Pairs{NewPair("3first", "value"),
	NewPair("3fourth", "value"), NewPair("3second", "value"), NewPair("3third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child4, err := NewLeafNode([]*Pairs{NewPair("4first", "value"),
	NewPair("4fourth", "value"), NewPair("4second", "value"), NewPair("4third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	parentNode.AddPoppedUpElementIntoCurrentNodeAndUpdateWithNewChildren(NewPair("popfirst","value"), child3, child4)

	child, err := parentNode.getChildAtIndex(0)
	if err != nil {
		t.Error(err)
	}
	if child.getElementAtIndex(0).key != "1first" {
		t.Error("Child not inserted at the correct position", child.getElements())
	}

	child, err = parentNode.getChildAtIndex(2)
	if err != nil {
		t.Error(err)
	}
	if child.getElementAtIndex(0).key !="4first" {
		printNodeElements(child)
		t.Error("Child not inserted at the correct position", child.getElements())
	}

}
