package main

import (
	"encoding/binary"
)

type IncomingPacket struct {
	buffer   []byte
	position int
}

func (p *IncomingPacket) init(size int) {
	p.buffer = make([]byte, size)
}

func (p *IncomingPacket) size() int {
	return len(p.buffer[p.position:])
}

func (p *IncomingPacket) resize(size int) {
	p.buffer = p.buffer[:size]
}

func (p *IncomingPacket) skipBytes(n int) {
	if p.position+n > p.size() {
		panic("skipping more bytes than size")
	}
	p.position += n
}

func (p *IncomingPacket) getUint8() uint8 {
	result := p.buffer[p.position]
	p.position += 1
	return result
}

func (p *IncomingPacket) peekUint8() uint8 {
	return p.buffer[p.position]
}

func (p *IncomingPacket) getUint16() uint16 {
	result := binary.LittleEndian.Uint16(p.buffer[p.position:])
	p.position += 2
	return result
}

func (p *IncomingPacket) peekUint16() uint16 {
	return binary.LittleEndian.Uint16(p.buffer[p.position:])
}

func (p *IncomingPacket) getUint32() uint32 {
	result := binary.LittleEndian.Uint32(p.buffer[p.position:])
	p.position += 4
	return result
}

func (p *IncomingPacket) peekUint32() uint32 {
	return binary.LittleEndian.Uint32(p.buffer[p.position:])
}

func (p *IncomingPacket) getString() string {
	stringLength := int(p.getUint16())
	result := string(p.buffer[p.position:(p.position + stringLength)])
	p.position += stringLength
	return result
}

func (p *IncomingPacket) peekBuffer() []byte {
	return p.buffer[p.position:]
}
