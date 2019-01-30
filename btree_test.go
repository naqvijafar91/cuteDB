package main

import (
	"fmt"
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
	for i := 1; i <= 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		if i == 230 {
			println("Inserted 229 elements")
		}
		tree.Insert(NewPair(key, value))
	}
	// tree.root.PrintTree()
}

func TestBtreeGet(t *testing.T) {
	tree, err := InitializeBtree(clearDB())
	if err != nil {
		t.Error(err)
	}
	totalElements := 500
	for i := 1; i <= totalElements; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		tree.Insert(NewPair(key, value))
	}

	for i := 1; i <= totalElements; i++ {
		key := fmt.Sprintf("key-%d", i)
		value, found, err := tree.Get(key)
		if err != nil {
			t.Error(err)
		}
		if !found || value == "" {
			t.Error("Value should be found ", key)
		}
	}

	for i := totalElements + 1; i <= totalElements+1+1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		_, found, err := tree.Get(key)
		if err != nil {
			t.Error(err)
		}
		if found {
			t.Error("Value should not be found")
		}
	}
}
