package cutedb

import "os"

type diskNodeService struct {
	file *os.File
}

func NewDiskNodeService(file *os.File) *diskNodeService {
	return &diskNodeService{file: file}
}
func (dns *diskNodeService) GetRootNodeFromDisk() (*DiskNode, error) {
	bs := NewBlockService(dns.file)
	rootBlock, err := bs.GetRootBlock()
	if err != nil {
		return nil, err
	}
	return bs.convertBlockToDiskNode(rootBlock), nil
}
