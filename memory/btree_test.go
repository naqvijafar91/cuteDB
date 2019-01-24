package memory

import "testing"

func TestBtreeInsert(t *testing.T) {
	tree := NewBtree()
	for i := 1; i <= 50; i++ {
		tree.Insert(int64(i))
	}
	// tree.root.PrintTree()
}
