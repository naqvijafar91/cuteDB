package cutedb

import "fmt"

/**
* Insertion Algorithm
1. It will begin from root and value will always be inserted into a leaf node
2. Insert Function Begin
3. If current node is leaf node, then return pick current node, Current Node Insertion Algorithm
    1. This section gives 3 outputs, 1 middle element and 2 child nodes or null,null,null
    2. Insert into the current node
    3. If its full, then sort it and make two new child nodes without the middle node ( NODE CREATION WILL TAKE PLACE HERE)
    4. take out the middle element along with the two child nodes,  Leaf Splitting no children Algorithm:
        1. Pick middle element by using length of array/2, lets say its index i
        2. Club all elements from 0 to i-1, and i+1 to len(array) and create new seperate nodes by inserting these 2 arrays into the respective keys[] of respective nodes
        3. Since the current node is a leaf node, we do not need to worry about its children and we can leave them to be null for both
        4. return middle,leftNode,rightNode
    5. If its not full, then return null,null,null
4. If this is not a leaf node, then find out the proper child node, Child Node Searching Algorithm:
    1. Input : Value to be inserted, the current Node. Output : Pointer to the childnode
    2. Since the list of values/elements is sorted, perform a binary or linear search to find the first element greater than the value to be inserted, if such an element is found, return pointer at position i, else return last pointer ( ie. the last pointer)
5. After getting the pointer to that element, call insert function Step 2 on that node RECURSIVELY ONLY HERE
6. If we get output from child node insert function Step 2, then this means that we have to insert the middle element received and accomodate the 2 pointers in the current node as well  discarding the old pointer ( NODE DESTRUCTION WILL ONLY TAKE PLACE HERE )
    1. If we got null as output then do nothing, else
    2. Insert into current Node, Popped up element and two child pointers insertion algorithm, Popped Up Joining Algorithm:
        1. Insert element and sort the array
        2. Now we need to discard 1 child pointer and insert 2 child pointers, Child Pointer Manipulation Algorithm :
        3. Find index of inserted element in array, lets say that it is i
        4. Now in the child pointer array, insert the left and right pointers at ith and i+1 th index
    3. If its full, sort it and make two new child nodes, Leaf Splitting with children Algorithm:
        1. Pick middle element by using length of array/2, lets say its index i (Same as 3.4.1)
        2. Club all elements from 0 to i-1, and i+1 to len(lkeys array) and create new seperate nodes by inserting these 2 arrays into the respective keys[] of respective nodes (Same as 3.4.2)
        3. For children[], split the current node's children array into 2 parts, part1 will be from 0 to i, and part 2 will be from i+1 to len(children array), and insert them into leftNode children, and rightNode children
        4. If current node is not the root node return middle,leftNode,rightNode
        5. else if current node == rootNode, Root Node Splitting Algorithm:
            1. Create a new node with elements array as keys[0] = middle
            2. children[0]=leftNode and children[1]=rightNode
            3. Set btree.root=new node
            4. return null,null,null

*/

const maxchildren = maxLeafSize + 1

// DiskNode - In memory node implementation
type DiskNode struct {
	keys             []*pairs
	childrenBlockIDs []uint64
	blockID          uint64
	blockService     *blockService
}

func (n *DiskNode) isLeaf() bool {
	if n.childrenBlockIDs == nil {
		return true
	} else if len(n.childrenBlockIDs) == 0 {
		return true
	}
	return false
}

// PrintTree - Traverse and print the entire tree
func (n *DiskNode) printTree(level ...int) {
	var currentLevel int
	if len(level) == 0 {
		currentLevel = 1
	} else {
		currentLevel = level[0]
	}
	n.printNode()
	for i := 0; i < len(n.childrenBlockIDs); i++ {
		fmt.Println("Printing ", i+1, " th child of level : ", currentLevel)
		childNode, err := n.getChildAtIndex(i)
		if err != nil {
			panic(err)
		}
		childNode.printTree(currentLevel + 1)
	}
}

/**
* Do a linear search and insert the element
 */
func (n *DiskNode) addElement(element *pairs) int {
	elements := n.getElements()
	indexForInsertion := 0
	elementInsertedInBetween := false
	for i := 0; i < len(elements); i++ {
		if elements[i].key >= element.key {
			// We have found the right place to insert the element

			indexForInsertion = i
			elements = append(elements, nil)
			copy(elements[indexForInsertion+1:], elements[indexForInsertion:])
			elements[indexForInsertion] = element
			n.setElements(elements)
			elementInsertedInBetween = true
			break
		}
	}
	if !elementInsertedInBetween {
		// If we are here, it means we need to insert the element at the rightmost position
		n.setElements(append(elements, element))
		indexForInsertion = len(n.getElements()) - 1
	}

	return indexForInsertion
}

func (n *DiskNode) hasOverFlown() bool {
	if len(n.getElements()) > n.blockService.getMaxLeafSize() {
		return true
	}
	return false
}

func (n *DiskNode) getElements() []*pairs {
	return n.keys
}

func (n *DiskNode) setElements(newElements []*pairs) {
	n.keys = newElements
}

func (n *DiskNode) getElementAtIndex(index int) *pairs {
	return n.keys[index]
}

func (n *DiskNode) getChildAtIndex(index int) (*DiskNode, error) {
	return n.blockService.getNodeAtBlockID(n.childrenBlockIDs[index])
}

func (n *DiskNode) shiftRemainingChildrenToRight(index int) {
	if len(n.childrenBlockIDs) < index+1 {
		// This means index is the last element, hence no need to shift
		return
	}
	n.childrenBlockIDs = append(n.childrenBlockIDs, 0)
	copy(n.childrenBlockIDs[index+1:], n.childrenBlockIDs[index:])
	n.childrenBlockIDs[index] = 0
}
func (n *DiskNode) setChildAtIndex(index int, childNode *DiskNode) {
	if len(n.childrenBlockIDs) < index+1 {
		n.childrenBlockIDs = append(n.childrenBlockIDs, 0)
	}
	n.childrenBlockIDs[index] = childNode.blockID
}

func (n *DiskNode) getLastChildNode() (*DiskNode, error) {
	return n.getChildAtIndex(len(n.childrenBlockIDs) - 1)
}

func (n *DiskNode) getChildNodes() ([]*DiskNode, error) {
	childNodes := make([]*DiskNode, len(n.childrenBlockIDs))
	for index := range n.childrenBlockIDs {
		childNode, err := n.getChildAtIndex(index)
		if err != nil {
			return nil, err
		}
		childNodes[index] = childNode
	}
	return childNodes, nil
}

func (n *DiskNode) getChildBlockIDs() []uint64 {
	return n.childrenBlockIDs
}

func (n *DiskNode) printNode() {
	fmt.Println("Printing Node")
	fmt.Println("--------------")
	for i := 0; i < len(n.getElements()); i++ {
		fmt.Println(n.getElementAtIndex(i))
	}
	fmt.Println("**********************")
}

// splitLeafNode - Split leaf node
func (n *DiskNode) splitLeafNode() (*pairs, *DiskNode, *DiskNode, error) {
	/**
		LEAF SPLITTING WITHOUT CHILDREN ALGORITHM
				If its full, then  make two new child nodes without the middle node ( NODE CREATION WILL TAKE PLACE HERE)
	    		Take out the middle element along with the two child nodes,  Leaf Splitting no children Algorithm:
	        	1. Pick middle element by using length of array/2, lets say its index i
	        	2. Club all elements from 0 to i-1, and i+1 to len(array) and create new seperate nodes by inserting these 2 arrays into the respective keys[] of respective nodes
	        	3. Since the current node is a leaf node, we do not need to worry about its children and we can leave them to be null for both
	        	4. return middle,leftNode,rightNode
	*/
	elements := n.getElements()
	midIndex := len(elements) / 2
	middle := elements[midIndex]

	// Now lets split elements array into 2 as we are splitting this node
	elements1 := elements[0:midIndex]
	elements2 := elements[midIndex+1 : len(elements)]

	// Now lets construct new Nodes from these 2 element arrays
	leftNode, err := newLeafNode(elements1, n.blockService)
	if err != nil {
		return nil, nil, nil, err
	}
	rightNode, err := newLeafNode(elements2, n.blockService)
	if err != nil {
		return nil, nil, nil, err
	}
	return middle, leftNode, rightNode, nil
}

//splitNonLeafNode - Split non leaf node
func (n *DiskNode) splitNonLeafNode() (*pairs, *DiskNode, *DiskNode, error) {
	/**
		NON-LEAF NODE SPLITTING ALGORITHM WITH CHILDREN MANIPULATION
		If its full, sort it and make two new child nodes, Leaf Splitting with children Algorithm:
	        1. Pick middle element by using length of array/2, lets say its index i (Same as 3.4.1)
			2. Club all elements from 0 to i-1, and i+1 to len(lkeys array) and create new seperate nodes
			   by inserting these 2 arrays into the respective keys[] of respective nodes (Same as 3.4.2)
			3. For children[], split the current node's children array into 2 parts, part1 will be
			   from 0 to i, and part 2 will be from i+1 to len(children array), and insert them into
			   leftNode children, and rightNode children

		NOTE : NODE CREATION WILL TAKE PLACE HERE
	*/
	elements := n.getElements()
	midIndex := len(elements) / 2
	middle := elements[midIndex]

	// Now lets split elements array into 2 as we are splitting this node
	elements1 := elements[0:midIndex]
	elements2 := elements[midIndex+1 : len(elements)]

	// Lets split the children
	children := n.childrenBlockIDs

	children1 := children[0 : midIndex+1]
	children2 := children[midIndex+1 : len(children)]

	// Now lets construct new Nodes from these 2 element arrays
	leftNode, err := newNodeWithChildren(elements1, children1, n.blockService)
	if err != nil {
		return nil, nil, nil, err
	}
	rightNode, err := newNodeWithChildren(elements2, children2, n.blockService)
	if err != nil {
		return nil, nil, nil, err
	}
	return middle, leftNode, rightNode, nil
}

// addPoppedUpElementIntoCurrentNodeAndUpdateWithNewChildren - Insert element received as a reaction
// from insert operation at child nodes
func (n *DiskNode) addPoppedUpElementIntoCurrentNodeAndUpdateWithNewChildren(element *pairs, leftNode *DiskNode, rightNode *DiskNode) {
	/**
		POPPED UP JOINING ALGORITHM
			Insert into current Node, Popped up element and two child pointers insertion algorithm, Popped Up Joining Algorithm:
	        1. Insert element and sort the array
	        2. Now we need to discard 1 child pointer and insert 2 child pointers, Child Pointer Manipulation Algorithm :
	        3. Find index of inserted element in array, lets say that it is i
	        4. Now in the child pointer array, insert the left and right pointers at ith and i+1 th index
	*/

	//CHILD POINTER MANIPULATION ALGORITHM
	insertionIndex := n.addElement(element)
	n.setChildAtIndex(insertionIndex, leftNode)
	//Shift remaining elements to the right and add this
	n.shiftRemainingChildrenToRight(insertionIndex + 1)
	n.setChildAtIndex(insertionIndex+1, rightNode)
}

// newLeafNode - Create a new leaf node without children
func newLeafNode(elements []*pairs, bs *blockService) (*DiskNode, error) {
	node := &DiskNode{keys: elements, blockService: bs}
	//persist the node to disk
	err := bs.saveNewNodeToDisk(node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// newNodeWithChildren - Create a non leaf node with children
func newNodeWithChildren(elements []*pairs, childrenBlockIDs []uint64, bs *blockService) (*DiskNode, error) {
	node := &DiskNode{keys: elements, childrenBlockIDs: childrenBlockIDs, blockService: bs}
	//persist this node to disk
	err := bs.saveNewNodeToDisk(node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// newRootNodeWithSingleElementAndTwoChildren - Create a new root node
func newRootNodeWithSingleElementAndTwoChildren(element *pairs, leftChildBlockID uint64,
	rightChildBlockID uint64, bs *blockService) (*DiskNode, error) {
	elements := []*pairs{element}
	childrenBlockIDs := []uint64{leftChildBlockID, rightChildBlockID}
	node := &DiskNode{keys: elements, childrenBlockIDs: childrenBlockIDs, blockService: bs}
	//persist this node to disk
	err := bs.updateRootNode(node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// getChildNodeForElement - Get Correct Traversal path for insertion
func (n *DiskNode) getChildNodeForElement(key string) (*DiskNode, error) {
	/** CHILD NODE SEARCHING ALGORITHM
		If this is not a leaf node, then find out the proper child node, Child Node Searching Algorithm:
	    1. Input : Value to be inserted, the current Node. Output : Pointer to the childnode
		2. Since the list of values/elements is sorted, perform a binary or linear search to find the
		   first element greater than the value to be inserted, if such an element is found, return pointer at position i, else return last pointer ( ie. the last pointer)
	*/

	for i := 0; i < len(n.getElements()); i++ {
		if key < n.getElementAtIndex(i).key {
			return n.getChildAtIndex(i)
		}
	}
	// This means that no element is found with value greater than the element to be inserted
	// so we need to return the last child node
	return n.getLastChildNode()
}

func (n *DiskNode) insert(value *pairs, bt *btree) (*pairs, *DiskNode, *DiskNode, error) {

	if n.isLeaf() {
		n.addElement(value)
		if !n.hasOverFlown() {
			// So lets store this updated node on disk
			err := n.blockService.updateNodeToDisk(n)
			if err != nil {
				return nil, nil, nil, err
			}
			return nil, nil, nil, nil
		}
		if bt.isRootNode(n) {
			poppedMiddleElement, leftNode, rightNode, err := n.splitLeafNode()
			if err != nil {
				return nil, nil, nil, err
			}
			//NOTE : NODE CREATION WILL TAKE PLACE HERE
			newRootNode, err := newRootNodeWithSingleElementAndTwoChildren(poppedMiddleElement,
				leftNode.blockID, rightNode.blockID, n.blockService)
			if err != nil {
				return nil, nil, nil, err
			}
			bt.setRootNode(newRootNode)
			return nil, nil, nil, nil

		}
		// Split the node and return to parent function with pooped up element and left,right nodes
		return n.splitLeafNode()

	}
	// Get the child Node for insertion
	childNodeToBeInserted, err := n.getChildNodeForElement(value.key)
	if err != nil {
		return nil, nil, nil, err
	}
	poppedMiddleElement, leftNode, rightNode, err := childNodeToBeInserted.insert(value, bt)
	if err != nil {
		return nil, nil, nil, err
	}
	if poppedMiddleElement == nil {
		// this means element has been inserted into the child and hence we do nothing
		return poppedMiddleElement, leftNode, rightNode, nil
	}
	// Insert popped up element into current node along with updating the child pointers
	// with new left and right nodes returned
	n.addPoppedUpElementIntoCurrentNodeAndUpdateWithNewChildren(poppedMiddleElement, leftNode, rightNode)

	if !n.hasOverFlown() {
		// this means that element has been easily inserted into current parent Node
		// without overflowing
		err := n.blockService.updateNodeToDisk(n)
		if err != nil {
			return nil, nil, nil, err
		}
		// So lets store this updated node on disk
		return nil, nil, nil, nil
	}
	// this means that the current parent node has overflown, we need to split this up
	// and move the popped up element upwards if this is not the root
	poppedMiddleElement, leftNode, rightNode, err = n.splitNonLeafNode()
	if err != nil {
		return nil, nil, nil, err
	}
	/**
		If current node is not the root node return middle,leftNode,rightNode
	    else if current node == rootNode, Root Node Splitting Algorithm:
	            1. Create a new node with elements array as keys[0] = middle
	            2. children[0]=leftNode and children[1]=rightNode
	            3. Set btree.root=new node
	            4. return null,null,null
	*/

	if !bt.isRootNode(n) {
		return poppedMiddleElement, leftNode, rightNode, nil
	}
	newRootNode, err := newRootNodeWithSingleElementAndTwoChildren(poppedMiddleElement,
		leftNode.blockID, rightNode.blockID, n.blockService)
	if err != nil {
		return nil, nil, nil, err
	}

	//@Todo: Update the metadata somewhere so that we can read this new root node
	//next time
	bt.setRootNode(newRootNode)
	return nil, nil, nil, nil
}

func (n *DiskNode) searchElementInNode(key string) (string, bool) {
	for i := 0; i < len(n.getElements()); i++ {
		if (n.getElementAtIndex(i)).key == key {
			return n.getElementAtIndex(i).value, true
		}
	}
	return "", false
}
func (n *DiskNode) search(key string) (string, error) {
	/*
		Algo:
		1. Find key in current node, if this is leaf node, then return as not found
		2. Then find the appropriate child node
		3. goto step 1
	*/
	value, foundInCurrentNode := n.searchElementInNode(key)

	if foundInCurrentNode {
		return value, nil
	}

	if n.isLeaf() {
		return "", nil
	}

	node, err := n.getChildNodeForElement(key)
	if err != nil {
		return "", err
	}
	return node.search(key)
}

// Insert - Insert value into Node
func (n *DiskNode) insertPair(value *pairs, bt *btree) error {
	_,_,_,err:=n.insert(value, bt)
	if err!=nil {
		return err
	}
	return nil
}

func (n *DiskNode) getValue(key string) (string, error) {
	return n.search(key)
}
