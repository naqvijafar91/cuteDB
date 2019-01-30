package cutedb

import "os"

type DiskNodeService struct {
	file *os.File
}

func NewDiskNodeService(file *os.File) *DiskNodeService {
	return &DiskNodeService{file: file}
}
func (dns *DiskNodeService) GetRootNodeFromDisk() (*DiskNode, error) {
	bs := NewBlockService(dns.file)
	rootBlock, err := bs.GetRootBlock()
	if err != nil {
		return nil, err
	}
	return bs.convertBlockToDiskNode(rootBlock),nil
}

