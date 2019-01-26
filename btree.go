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

//@Todo: Remove panic from here
// NewBtree - Create a new btree
func InitializeBtree() (*Btree, error) {
	path := "./db/freedom.db"
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	dns := NewDiskNodeService(file)

	root, err := dns.GetRootNodeFromDisk()
	if err != nil {
		panic(err)
	}
	return &Btree{root: root}, nil
}

// Insert - Insert element in tree
func (btree *Btree) Insert(value int64) {
	btree.root.Insert(value, btree)
}

func (btree *Btree) setRootNode(node Node) {
	btree.root = node
}
