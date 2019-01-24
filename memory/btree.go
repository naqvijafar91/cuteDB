package memory

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

// NewBtree - Create a new btree
func NewBtree() *Btree {
	return &Btree{root: NewLeafNode([]int64{})}
}

// Insert - Insert element in tree
func (btree *Btree) Insert(value int64) {
	btree.root.Insert(value, btree)
}

func (btree *Btree) setRootNode(node Node) {
	btree.root = node
}
