package main

import (
	"os"
	"testing"
)

func clearDB() string {
	path := "./db/test.db"
	if _, err := os.Stat(path); err == nil {
		// path/to/whatever exists
		err := os.Remove(path)
		if err != nil {
			panic(err)
		}
	}
	return path
}

func TestBtreeInsert(t *testing.T) {
	tree, err := InitializeBtree(clearDB())
	if err != nil {
		t.Error(err)
	}
	for i := 1; i <= 50; i++ {
		tree.Insert(int64(i))
	}
	// tree.root.PrintTree()
}

func TestBtreeGet(t *testing.T) {
	tree, err := InitializeBtree(clearDB())
	if err != nil {
		t.Error(err)
	}
	totalElements := 5000
	for i := 1; i <= totalElements; i++ {
		tree.Insert(int64(i))
	}

	for i := 1; i <= totalElements; i++ {
		found, err := tree.Get(int64(i))
		if err != nil {
			t.Error(err)
		}
		if !found {
			t.Error("Value should be found")
		}
	}

	for i := totalElements + 1; i <= totalElements+1+1000; i++ {
		found, err := tree.root.Get(int64(i))
		if err != nil {
			t.Error(err)
		}
		if found {
			t.Error("Value should not be found")
		}
	}
}
