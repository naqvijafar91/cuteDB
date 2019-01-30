package cutedb

import "os"

type diskNodeService struct {
	file *os.File
}

func newDiskNodeService(file *os.File) *diskNodeService {
	return &diskNodeService{file: file}
}
func (dns *diskNodeService) getRootNodeFromDisk() (*DiskNode, error) {
	bs := newBlockService(dns.file)
	rootBlock, err := bs.getRootBlock()
	if err != nil {
		return nil, err
	}
	return bs.convertBlockToDiskNode(rootBlock), nil
}
