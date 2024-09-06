package main

import (
	"encoding/binary"
	"fmt"
)

const (
	HEADER_OFFSET = 10
)

type OutgoingPacket struct {
	buffer   []byte
	position int
	header   int
}

func (p *OutgoingPacket) init(size int) {
	p.buffer = make([]byte, size)
	p.header = 10
	p.position = 0
}

func (p *OutgoingPacket) size() int {
	return p.position + (HEADER_OFFSET - p.header)
}

func (p *OutgoingPacket) get() []byte {
	return p.buffer[p.header:(HEADER_OFFSET + p.position)]
}

func (p *OutgoingPacket) addUint8(data uint8) {
	offset := HEADER_OFFSET + p.position
	if (offset + 1) > len(p.buffer) {
		fmt.Println("Error: Buffer overflow")
		return
	}
	p.buffer[offset] = data
	p.position += 1
}

func (p *OutgoingPacket) addUint16(data uint16) {
	offset := HEADER_OFFSET + p.position
	if (offset + 2) > len(p.buffer) {
		fmt.Println("Error: Buffer overflow")
		return
	}
	binary.LittleEndian.PutUint16(p.buffer[offset:], data)
	p.position += 2
}

func (p *OutgoingPacket) addUint32(data uint32) {
	offset := HEADER_OFFSET + p.position
	if (offset + 4) > len(p.buffer) {
		fmt.Println("Error: Buffer overflow")
		return
	}
	binary.LittleEndian.PutUint32(p.buffer[offset:], data)
	p.position += 4
}

func (p *OutgoingPacket) addString(data string) {
	stringLength := len(data)
	if stringLength > 65535 { // Maximum size of a uint16
		fmt.Println("Error: String is too long")
		return
	}
	p.addUint16(uint16(stringLength))

	offset := HEADER_OFFSET + p.position
	if offset+stringLength > len(p.buffer) {
		fmt.Println("Error: Buffer overflow")
		return
	}

	copy(p.buffer[offset:], []byte(data))
	p.position += stringLength
}

func (p *OutgoingPacket) headerAddSize() {
	size := uint16(p.size())
	binary.LittleEndian.PutUint16(p.buffer[p.header-2:], size)
	p.header -= 2
}

func (p *OutgoingPacket) addPadding() {
	size := p.size()
	if size%8 != 0 {
		toAdd := 8 - (size % 8)
		for i := 0; i < toAdd; i++ {
			p.addUint8(0x33)
		}
	}
}

func (p *OutgoingPacket) xteaEncrypt(xteaKey [4]uint32) error {

	p.headerAddSize()
	p.addPadding()

	expandedXteaKey := expandXteaKey(xteaKey)
	xteaEncrypt(p.buffer[p.header:], expandedXteaKey)

	return nil
}
