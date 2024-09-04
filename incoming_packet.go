package main

import (
	"encoding/binary"
)

type IncomingPacket struct {
	buffer []byte
}

// Use pointer receiver to modify the buffer
func (p *IncomingPacket) init(size int) {
	p.buffer = make([]byte, size)
}

func (p *IncomingPacket) size() int {
	return len(p.buffer)
}

// Use pointer receiver to resize the buffer
func (p *IncomingPacket) resize(size int) {
	p.buffer = p.buffer[:size]
}

// Skip n bytes in the buffer
func (p *IncomingPacket) skipBytes(n int) {
	p.buffer = p.buffer[n:]
}

// Get a uint8 from the buffer
func (p *IncomingPacket) getUint8() uint8 {
	result := p.buffer[0]
	p.buffer = p.buffer[1:]
	return result
}

// Get a uint16 from the buffer
func (p *IncomingPacket) getUint16() uint16 {
	result := binary.LittleEndian.Uint16(p.buffer)
	p.buffer = p.buffer[2:]
	return result
}

// Get a uint32 from the buffer
func (p *IncomingPacket) getUint32() uint32 {
	result := binary.LittleEndian.Uint32(p.buffer)
	p.buffer = p.buffer[4:]
	return result
}

// Get a string from the buffer
func (p *IncomingPacket) getString() string {
	stringLength := p.getUint16()
	result := string(p.buffer[:stringLength])
	p.buffer = p.buffer[stringLength:]
	return result
}
