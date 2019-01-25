package main

import "os"
// Btree - Our in memory Btree struct
type Btree struct {
	root Node
}

// Node - Interface for node
type Node interface {
	Insert(value int64, btree *Btree)
	PrintTree(level ...int)
}

func (btree *Btree) isRootNode(node Node) bool {
	return btree.root == node
}

func InitDb() *Btree {
	return &Btree{root: InitRootNode()}
}

// NewBtree - Create a new btree
func NewBtree() *Btree {
	path := "./db/freedom.db"
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	blockService:=NewBlockService(file)
	return &Btree{root: NewLeafNode([]int64{},blockService)}
}

// Insert - Insert element in tree
func (btree *Btree) Insert(value int64) {
	btree.root.Insert(value, btree)
}

func (btree *Btree) setRootNode(node Node) {
	btree.root = node
}
