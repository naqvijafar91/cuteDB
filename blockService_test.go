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

func TestShouldConvertToAndFromDiskNode(t *testing.T) {
	bs := initBlockService()
	node := &DiskNode{}
	node.blockID = 55
	node.keys = []int64{500, 100}
	node.childrenBlockIDs=[]uint64{1000,10001}
	block:=bs.convertDiskNodeToBlock(node)
	
	if block.id != 55 {
		t.Error("Should have same block id as node block id")
	}
	if block.data[1] != 100 {
		t.Error("Should have same data element as node")
	}

	if block.childrenBlockIds[1]!=10001 {
		t.Error("Block ids should match")
	}

	nodeFromBlock:=bs.convertBlockToDiskNode(block)
	
	if nodeFromBlock.blockID!=node.blockID {
		t.Error("Block ids should match")
	}

	if nodeFromBlock.childrenBlockIDs[0] !=1000 {
		t.Error("Child Block ids should match")
	}
	if nodeFromBlock.keys[0] != 500 {
		t.Error("Data elements should match")
	}
}
