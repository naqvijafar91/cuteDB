package main

import "os"

// Btree - Our in memory Btree struct
type Btree struct {
	root Node
}

// Node - Interface for node
type Node interface {
	Insert(value *Pairs,btree *Btree)
	Get(key string) (bool, error)
	PrintTree(level ...int)
}

func (btree *Btree) isRootNode(node Node) bool {
	return btree.root == node
}

//@Todo: Remove panic from here
// NewBtree - Create a new btree
func InitializeBtree(path ...string) (*Btree, error) {
	if len(path) == 0 {
		path = make([]string, 1)
		path[0] = "./db/freedom.db"
	}

	file, err := os.OpenFile(path[0], os.O_RDWR|os.O_CREATE, 0666)
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
func (btree *Btree) Insert(value *Pairs) {
	btree.root.Insert(value, btree)
}

func (btree *Btree) Get(key string) (bool,error) {
	return btree.root.Get(key)
}
func (btree *Btree) setRootNode(node Node) {
	btree.root = node
}
