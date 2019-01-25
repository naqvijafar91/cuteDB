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

func TestShouldConvertToAndFromBytes(t *testing.T) {
	blockService := initBlockService()
	block := &Block{}
	block.setData([]uint64{100, 101, 102})
	block.setChildren([]uint64{2, 3, 4, 6})
	blockBuffer := blockService.getBufferFromBlock(block)
	convertedBlock := blockService.getBlockFromBuffer(blockBuffer)

	if len(convertedBlock.data) != len(block.data) {
		t.Error("Length of blocks should be same")
	}

	if convertedBlock.data[1] != 101 {
		t.Error("Should contain 101 at 1st index")
	}

	if convertedBlock.childrenBlockIds[2] != 4 {
		t.Error("Should contain 4 at 2nd index")
	}

}
