package cutedb

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
	elements := make([]*pairs, 3)
	elements[0] = newPair("hola", "amigos")
	elements[1] = newPair("foo", "bar")
	elements[2] = newPair("gooz", "bumps")
	n, err := newLeafNode(elements, blockService)
	if err != nil {
		t.Error(err)
	}
	addedElement := newPair("added", "please check")
	n.addElement(addedElement)

	if !reflect.DeepEqual(n.getElements(), []*pairs{addedElement, elements[0],
		elements[1], elements[2]}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}

	n, err = newLeafNode([]*pairs{newPair("first", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	n.addElement(newPair("second", "value"))
	if !reflect.DeepEqual(n.getElements(), []*pairs{newPair("first", "value"),
		newPair("second", "value")}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}

	n, err = newLeafNode([]*pairs{newPair("first", "value"),
		newPair("second", "value"), newPair("third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	n.addElement(newPair("fourth", "value"))
	if !reflect.DeepEqual(n.getElements(), []*pairs{newPair("first", "value"),
		newPair("fourth", "value"), newPair("second", "value"), newPair("third", "value")}) {
		t.Error("Value not inserted at the correct position", n.getElements())
	}
}

func TestIsLeaf(t *testing.T) {
	blockService := initBlockService()
	child1, err := newLeafNode([]*pairs{newPair("first", "value"),
		newPair("second", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child2, err := newLeafNode([]*pairs{newPair("third", "value"),
		newPair("forth", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	n, err := newNodeWithChildren([]*pairs{newPair("fifth", "value"),
		newPair("sixth", "value")}, []uint64{child1.blockID, child2.blockID}, blockService)
	if err != nil {
		t.Error(err)
	}
	if n.isLeaf() {
		t.Error("Should not return as leaf as it has children", n)
	}

	child1, err = newLeafNode(nil, blockService)
	if err != nil {
		t.Error(err)
	}
	child2, err = newLeafNode(nil, blockService)
	if err != nil {
		t.Error(err)
	}
	n, err = newNodeWithChildren([]*pairs{newPair("first", "value"),
		newPair("second", "value")}, nil, blockService)
	if err != nil {
		t.Error(err)
	}
	if !n.isLeaf() {
		t.Error("Should return as leaf as it has no children", n)
	}

	n, err = newLeafNode([]*pairs{newPair("first", "value"),
		newPair("second", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	if !n.isLeaf() {
		t.Error("Should return as leaf as it has no children", n)
	}
}

func TestHasOverFlown(t *testing.T) {
	blockService := initBlockService()
	elements := make([]*pairs, blockService.getMaxLeafSize()+1)
	for i := 0; i < blockService.getMaxLeafSize()+1; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		elements[i] = newPair(key, value)
	}
	n, err := newLeafNode(elements, blockService)
	if err != nil {
		t.Error(err)
	}
	if !n.hasOverFlown() {
		t.Error("Should return true as node has overflown", n)
	}

	n, err = newLeafNode([]*pairs{newPair("first", "value"), newPair("fourth", "value"),
		newPair("second", "value"), newPair("third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	if n.hasOverFlown() {
		t.Error("Should return false as node has 3 elements", n)
	}

	child1, err := newLeafNode([]*pairs{newPair("first", "value"),
		newPair("second", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child2, err := newLeafNode([]*pairs{newPair("third", "value"),
		newPair("fourth", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	n, err = newNodeWithChildren(elements, []uint64{child1.blockID,
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
	n, err := newLeafNode([]*pairs{newPair("first", "value"),
		newPair("fourth", "value"), newPair("second", "value"), newPair("third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	poppedUpMiddleElement, leftChild, rightChild, err := n.splitLeafNode()
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
	child1, err := newLeafNode([]*pairs{newPair("1first", "value"),
		newPair("1fourth", "value"), newPair("1second", "value"), newPair("1third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child2, err := newLeafNode([]*pairs{newPair("2first", "value"),
		newPair("2fourth", "value"), newPair("2second", "value"), newPair("2third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child3, err := newLeafNode([]*pairs{newPair("3first", "value"),
		newPair("3fourth", "value"), newPair("3second", "value"), newPair("3third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child4, err := newLeafNode([]*pairs{newPair("4first", "value"),
		newPair("4fourth", "value"), newPair("4second", "value"), newPair("4third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child5, err := newLeafNode([]*pairs{newPair("5first", "value"),
		newPair("5fourth", "value"), newPair("5second", "value"), newPair("5third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}

	n, err := newNodeWithChildren([]*pairs{newPair("nfirst", "value"),
		newPair("nfourth", "value"), newPair("nsecond", "value"), newPair("nthird", "value")},
		[]uint64{child1.blockID, child2.blockID, child3.blockID,
			child4.blockID, child5.blockID}, blockService)
	if err != nil {
		t.Error(err)
	}
	poppedUpMiddleElement, leftChild, rightChild, err := n.splitNonLeafNode()
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
	child1OfParent, err := newLeafNode([]*pairs{newPair("1first", "value"),
		newPair("1fourth", "value"), newPair("1second", "value"), newPair("1third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child2OfParent, err := newLeafNode([]*pairs{newPair("2first", "value"),
		newPair("2fourth", "value"), newPair("2second", "value"), newPair("2third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	parentNode, err := newNodeWithChildren([]*pairs{newPair("parentfirst", "value")}, []uint64{child1OfParent.blockID,
		child2OfParent.blockID}, blockService)
	child3, err := newLeafNode([]*pairs{newPair("3first", "value"),
		newPair("3fourth", "value"), newPair("3second", "value"), newPair("3third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	child4, err := newLeafNode([]*pairs{newPair("4first", "value"),
		newPair("4fourth", "value"), newPair("4second", "value"), newPair("4third", "value")}, blockService)
	if err != nil {
		t.Error(err)
	}
	parentNode.addPoppedUpElementIntoCurrentNodeAndUpdateWithNewChildren(newPair("popfirst", "value"), child3, child4)

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
	if child.getElementAtIndex(0).key != "4first" {
		printNodeElements(child)
		t.Error("Child not inserted at the correct position", child.getElements())
	}

}
