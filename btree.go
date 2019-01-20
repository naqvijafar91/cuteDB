package main

// Btree - Our in memory Btree struct
type Btree struct {
	root Node
}

type Node interface {
	Insert(value int64, btree *Btree)
	PrintTree(level ...int)
}

func (btree *Btree) isRootNode(node Node) bool {
	return btree.root == node
}

func NewBtree() *Btree {
	return &Btree{root: NewLeafNode([]int64{})}
}
func (tree *Btree) Insert(value int64) {
	tree.root.Insert(value, tree)
}
