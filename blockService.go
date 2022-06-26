package cutedb

import (
	"encoding/binary"
	"os"
)

const blockSize = 4096

// Based on the below calc
const maxLeafSize = 30

// diskBlock -- Make sure that it is accomodated in blockSize = 4096
type diskBlock struct {
	id                  uint64   // 4096 - 8 = 4088
	currentLeafSize     uint64   // 4088 - 8 = 4080
	currentChildrenSize uint64   // 4080 - 8 = 4072
	childrenBlockIds    []uint64 // 262 - (8 * 30) =  22
	dataSet             []*pairs // 4072 - (127 * 30) = 262
}

// 22 bytes are still wasted

func (b *diskBlock) setData(data []*pairs) {
	b.dataSet = data
	b.currentLeafSize = uint64(len(data))
}

func (b *diskBlock) setChildren(childrenBlockIds []uint64) {
	b.childrenBlockIds = childrenBlockIds
	b.currentChildrenSize = uint64(len(childrenBlockIds))
}

func uint64ToBytes(index uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(index))
	return b
}

func uint64FromBytes(b []byte) uint64 {
	return uint64(binary.LittleEndian.Uint64(b))
}

type blockService struct {
	file *os.File
	mu *sync.Mutex
}

func (bs *blockService) getLatestBlockID() (int64, error) {

	fi, err := bs.file.Stat()
	if err != nil {
		return -1, err
	}

	length := fi.Size()
	if length == 0 {
		return -1, nil
	}
	// Calculate page number required to be fetched from disk
	return (int64(fi.Size()) / int64(blockSize)) - 1, nil
}

//@Todo:Store current root block data somewhere else
func (bs *blockService) getRootBlock() (*diskBlock, error) {

	/*
		1. Check if root block exists
		2. If exisits, fetch it, else initialize a new block
	*/
	if !bs.rootBlockExists() {
		// Need to write a new block
		return bs.newBlock()

	}
	return bs.getBlockFromDiskByBlockNumber(0)

}

func (bs *blockService) getBlockFromDiskByBlockNumber(index int64) (*diskBlock, error) {
	if index < 0 {
		panic("Index less than 0 asked")
	}
	offset := index * blockSize
	_, err := bs.file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	blockBuffer := make([]byte, blockSize)
	_, err = bs.file.Read(blockBuffer)
	if err != nil {
		return nil, err
	}
	block := bs.getBlockFromBuffer(blockBuffer)
	return block, nil
}

func (bs *blockService) getBlockFromBuffer(blockBuffer []byte) *diskBlock {
	blockOffset := 0
	block := &diskBlock{}

	//Read Block index
	block.id = uint64FromBytes(blockBuffer[blockOffset:])
	blockOffset += 8
	block.currentLeafSize = uint64FromBytes(blockBuffer[blockOffset:])
	blockOffset += 8
	block.currentChildrenSize = uint64FromBytes(blockBuffer[blockOffset:])
	blockOffset += 8
	//Read actual pairs now
	block.dataSet = make([]*pairs, block.currentLeafSize)
	for i := 0; i < int(block.currentLeafSize); i++ {
		block.dataSet[i] = convertBytesToPair(blockBuffer[blockOffset:])
		blockOffset += pairSize
	}
	// Read children block indexes
	block.childrenBlockIds = make([]uint64, block.currentChildrenSize)
	for i := 0; i < int(block.currentChildrenSize); i++ {
		block.childrenBlockIds[i] = uint64FromBytes(blockBuffer[blockOffset:])
		blockOffset += 8
	}
	return block
}

func (bs *blockService) getBufferFromBlock(block *diskBlock) []byte {
	blockBuffer := make([]byte, blockSize)
	blockOffset := 0

	//Write Block index
	copy(blockBuffer[blockOffset:], uint64ToBytes(block.id))
	blockOffset += 8
	copy(blockBuffer[blockOffset:], uint64ToBytes(block.currentLeafSize))
	blockOffset += 8
	copy(blockBuffer[blockOffset:], uint64ToBytes(block.currentChildrenSize))
	blockOffset += 8

	//Write actual pairs now
	for i := 0; i < int(block.currentLeafSize); i++ {
		copy(blockBuffer[blockOffset:], convertPairsToBytes(block.dataSet[i]))
		blockOffset += pairSize
	}
	// Read children block indexes
	for i := 0; i < int(block.currentChildrenSize); i++ {
		copy(blockBuffer[blockOffset:], uint64ToBytes(block.childrenBlockIds[i]))
		blockOffset += 8
	}
	return blockBuffer
}

func (bs *blockService) newBlock() (*diskBlock, error) {

	latestBlockID, err := bs.getLatestBlockID()
	block := &diskBlock{}
	if err != nil {
		// This means that no file exists
		block.id = 0
	} else {
		block.id = uint64(latestBlockID) + 1
	}
	block.currentLeafSize = 0
	err = bs.writeBlockToDisk(block)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (bs *blockService) writeBlockToDisk(block *diskBlock) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	seekOffset := blockSize * block.id
	blockBuffer := bs.getBufferFromBlock(block)
	_, err := bs.file.Seek(int64(seekOffset), 0)
	if err != nil {
		return err
	}
	_, err = bs.file.Write(blockBuffer)
	if err != nil {
		return err
	}
	return nil
}

func (bs *blockService) convertDiskNodeToBlock(node *DiskNode) *diskBlock {
	block := &diskBlock{id: node.blockID}
	tempElements := make([]*pairs, len(node.getElements()))
	for index, element := range node.getElements() {
		tempElements[index] = element
	}
	block.setData(tempElements)
	tempBlockIDs := make([]uint64, len(node.getChildBlockIDs()))
	for index, childBlockID := range node.getChildBlockIDs() {
		tempBlockIDs[index] = childBlockID
	}
	block.setChildren(tempBlockIDs)
	return block
}

func (bs *blockService) getNodeAtBlockID(blockID uint64) (*DiskNode, error) {
	block, err := bs.getBlockFromDiskByBlockNumber(int64(blockID))
	if err != nil {
		return nil, err
	}
	return bs.convertBlockToDiskNode(block), nil
}

func (bs *blockService) convertBlockToDiskNode(block *diskBlock) *DiskNode {
	node := &DiskNode{
		blockID:      block.id,
		blockService: bs,
		keys:         make([]*pairs, block.currentLeafSize),
	}
	for index := range node.keys {
		node.keys[index] = block.dataSet[index]
	}
	node.childrenBlockIDs = make([]uint64, block.currentChildrenSize)
	for index := range node.childrenBlockIDs {
		node.childrenBlockIDs[index] = block.childrenBlockIds[index]
	}
	return node
}

// NewBlockFromNode - Save a new node to disk block
func (bs *blockService) saveNewNodeToDisk(n *DiskNode) error {
	// Get block id to be assigned to this block
	latestBlockID, err := bs.getLatestBlockID()
	if err != nil {
		return err
	}
	n.blockID = uint64(latestBlockID) + 1
	block := bs.convertDiskNodeToBlock(n)
	return bs.writeBlockToDisk(block)
}

func (bs *blockService) updateNodeToDisk(n *DiskNode) error {
	block := bs.convertDiskNodeToBlock(n)
	return bs.writeBlockToDisk(block)
}

func (bs *blockService) updateRootNode(n *DiskNode) error {
	n.blockID = 0
	return bs.updateNodeToDisk(n)
}

func newBlockService(file *os.File) *blockService {
	return &blockService{file: file, mu: &sync.Mutex{}}
}

func (bs *blockService) rootBlockExists() bool {
	latestBlockID, err := bs.getLatestBlockID()
	// fmt.Println(latestBlockID)
	//@Todo:Validate the type of error here
	if err != nil {
		// Need to write a new block
		return false
	} else if latestBlockID == -1 {
		return false
	} else {
		return true
	}
}

/**
@Todo: Implement a function to :
1. Dynamicaly calculate blockSize
2. Then based on the blocksize, calculate the maxLeafSize
*/
func (bs *blockService) getMaxLeafSize() int {
	return maxLeafSize
}
