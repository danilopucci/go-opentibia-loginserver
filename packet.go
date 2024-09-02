package main

import (
	"encoding/binary"
)

type Packet struct {
	buffer []byte
}

// Use pointer receiver to modify the buffer
func (p *Packet) init(size int) {
	p.buffer = make([]byte, size)
}

// Use pointer receiver to resize the buffer
func (p *Packet) resize(size int) {
	p.buffer = p.buffer[:size]
}

// Skip n bytes in the buffer
func (p *Packet) skipBytes(n int) {
	p.buffer = p.buffer[n:]
}

// Get a uint8 from the buffer
func (p *Packet) getUint8() uint8 {
	result := p.buffer[0]
	p.buffer = p.buffer[1:]
	return result
}

// Get a uint16 from the buffer
func (p *Packet) getUint16() uint16 {
	result := binary.LittleEndian.Uint16(p.buffer)
	p.buffer = p.buffer[2:]
	return result
}

// Get a uint32 from the buffer
func (p *Packet) getUint32() uint32 {
	result := binary.LittleEndian.Uint32(p.buffer)
	p.buffer = p.buffer[4:]
	return result
}

// Get a string from the buffer
func (p *Packet) getString() string {
	stringLength := p.getUint16()
	result := string(p.buffer[:stringLength])
	p.buffer = p.buffer[stringLength:]
	return result
}