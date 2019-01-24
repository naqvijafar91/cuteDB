package main

import (
	"os"
	"testing"
)

func initBlockService() *BlockService {
	path := "./db/test.db"
	if _, err := os.Stat(path); err == nil {
		// path/to/whatever exists
		err := os.Remove(path)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	return NewBlockService(file)
}

func TestShouldGetNegativeIfBlockNotPresent(t *testing.T) {
	blockService := initBlockService()
	latestBlockID, _ := blockService.getLatestBlockID()
	if latestBlockID != -1 {
		t.Error("Should get negative block id")
	}
}

func TestShouldSuccessfullyInitializeNewBlock(t *testing.T) {
	blockService := initBlockService()
	block, err := blockService.GetRootBlock()
	if err != nil {
		t.Error(err)
	}
	if block.id != 0 {
		t.Error("Root Block id should be zero")
	}

	if block.currentLeafSize != 0 {
		t.Error("Block leaf size should be zero")
	}
}

func TestShouldSaveNewBlockOnDisk(t *testing.T) {
	blockService := initBlockService()
	block, err := blockService.GetRootBlock()
	if err != nil {
		t.Error(err)
	}
	if block.id != 0 {
		t.Error("Root Block id should be zero")
	}

	if block.currentLeafSize != 0 {
		t.Error("Block leaf size should be zero")
	}

	block.setData([]uint64{55, 100})
	err = blockService.writeBlockToDisk(block)
	if err != nil {
		t.Error(err)
	}

	block, err = blockService.GetRootBlock()
	if err != nil {
		t.Error(err)
	}

	if len(block.data) == 0 {
		t.Error("Length of data field should not be zero")
	}
}
