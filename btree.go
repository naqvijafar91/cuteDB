package cutedb

import "os"

// Btree - Our in memory Btree struct
type Btree struct {
	root node
}

// node - Interface for node
type node interface {
	Insert(value *Pairs, btree *Btree)
	Get(key string) (string, error)
	PrintTree(level ...int)
}

func (btree *Btree) isRootNode(n node) bool {
	return btree.root == n
}

// NewBtree - Create a new btree
func InitializeBtree(path ...string) (*Btree, error) {
	if len(path) == 0 {
		path = make([]string, 1)
		path[0] = "./db/freedom.db"
	}

	file, err := os.OpenFile(path[0], os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil,err
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

func (btree *Btree) Get(key string) (string, bool, error) {
	value, err := btree.root.Get(key)
	if err != nil {
		return "", false, err
	}
	if value == "" {
		return "", false, nil
	}
	return value, true, nil
}

func (btree *Btree) setRootNode(n node) {
	btree.root = n
}
