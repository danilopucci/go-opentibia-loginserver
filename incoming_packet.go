package main

import (
	"encoding/binary"
)

type IncomingPacket struct {
	buffer   []byte
	position int
}

// Use pointer receiver to modify the buffer
func (p *IncomingPacket) init(size int) {
	p.buffer = make([]byte, size)
}

func (p *IncomingPacket) size() int {
	return len(p.buffer[p.position:])
}

// Use pointer receiver to resize the buffer
func (p *IncomingPacket) resize(size int) {
	p.buffer = p.buffer[:size]
}

// Skip n bytes in the buffer
func (p *IncomingPacket) skipBytes(n int) {
	if p.position+n > p.size() {
		panic("skipping more bytes than size")
	}
	p.position += n
}

// Get a uint8 from the buffer
func (p *IncomingPacket) getUint8() uint8 {
	result := p.buffer[p.position]
	p.position += 1
	return result
}

// Get a uint16 from the buffer
func (p *IncomingPacket) getUint16() uint16 {
	result := binary.LittleEndian.Uint16(p.buffer[p.position:])
	p.position += 2
	return result
}

// Get a uint32 from the buffer
func (p *IncomingPacket) getUint32() uint32 {
	result := binary.LittleEndian.Uint32(p.buffer[p.position:])
	p.position += 4
	return result
}

// Get a string from the buffer
func (p *IncomingPacket) getString() string {
	stringLength := int(p.getUint16())
	result := string(p.buffer[p.position:(p.position + stringLength)])
	p.position += stringLength
	return result
}
