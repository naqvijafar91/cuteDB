package cuteDB

import "encoding/binary"

// 2+2+30+90 = 124
const pairSize = 124

type Pairs struct {
	keyLen   uint16 // 2
	valueLen uint16 // 2
	key      string // 30
	value    string // 90
}

func (p *Pairs) setKey(key string) {
	p.key = key
	p.keyLen = uint16(len(key))
}

func (p *Pairs) setValue(value string) {
	p.value = value
	p.valueLen = uint16(len(value))
}

func NewPair(key string, value string) *Pairs {
	pair := &Pairs{}
	pair.setKey(key)
	pair.setValue(value)
	return pair
}

func convertPairsToBytes(pair *Pairs) []byte {
	pairByte := make([]byte, pairSize)
	var pairOffset uint16
	pairOffset = 0
	copy(pairByte[pairOffset:], uint16ToBytes(pair.keyLen))
	pairOffset += 2
	copy(pairByte[pairOffset:], uint16ToBytes(pair.valueLen))
	pairOffset += 2
	keyByte := []byte(pair.key)
	copy(pairByte[pairOffset:], keyByte[:pair.keyLen])
	pairOffset += pair.keyLen
	valueByte := []byte(pair.value)
	copy(pairByte[pairOffset:], valueByte[:pair.valueLen])
	return pairByte
}

func convertBytesToPair(pairByte []byte) *Pairs {
	pair := &Pairs{}
	var pairOffset uint16
	pairOffset = 0
	//Read key length
	pair.keyLen = uint16FromBytes(pairByte[pairOffset:])
	pairOffset += 2
	//Read value length
	pair.valueLen = uint16FromBytes(pairByte[pairOffset:])
	pairOffset += 2
	pair.key = string(pairByte[pairOffset : pairOffset+pair.keyLen])
	pairOffset += pair.keyLen
	pair.value = string(pairByte[pairOffset : pairOffset+pair.valueLen])
	return pair
}

func uint16FromBytes(b []byte) uint16 {
	i := uint16(binary.LittleEndian.Uint64(b))
	return i
}

func uint16ToBytes(value uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(value))
	return b
}
