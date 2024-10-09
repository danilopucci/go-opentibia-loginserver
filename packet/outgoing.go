package packet

import (
	"encoding/binary"
	"fmt"
	"go-opentibia-loginserver/crypt"
)

const (
	HEADER_OFFSET = 10
)

type Outgoing struct {
	buffer   []byte
	position int
	header   int
}

func NewOutgoing(size int) *Outgoing {
	return &Outgoing{
		buffer:   make([]byte, size),
		position: 0,
		header:   HEADER_OFFSET,
	}
}

func (p *Outgoing) size() int {
	return p.position + (HEADER_OFFSET - p.header)
}

func (p *Outgoing) Get() []byte {
	return p.buffer[p.header:(HEADER_OFFSET + p.position)]
}

func (p *Outgoing) AddUint8(data uint8) {
	offset := HEADER_OFFSET + p.position
	if (offset + 1) > len(p.buffer) {
		fmt.Println("Error: Buffer overflow")
		return
	}
	p.buffer[offset] = data
	p.position += 1
}

func (p *Outgoing) AddUint16(data uint16) {
	offset := HEADER_OFFSET + p.position
	if (offset + 2) > len(p.buffer) {
		fmt.Println("Error: Buffer overflow")
		return
	}
	binary.LittleEndian.PutUint16(p.buffer[offset:], data)
	p.position += 2
}

func (p *Outgoing) AddUint32(data uint32) {
	offset := HEADER_OFFSET + p.position
	if (offset + 4) > len(p.buffer) {
		fmt.Println("Error: Buffer overflow")
		return
	}
	binary.LittleEndian.PutUint32(p.buffer[offset:], data)
	p.position += 4
}

func (p *Outgoing) AddString(data string) {
	stringLength := len(data)
	if stringLength > 65535 { // Maximum size of a uint16
		fmt.Println("Error: String is too long")
		return
	}
	p.AddUint16(uint16(stringLength))

	offset := HEADER_OFFSET + p.position
	if offset+stringLength > len(p.buffer) {
		fmt.Println("Error: Buffer overflow")
		return
	}

	copy(p.buffer[offset:], []byte(data))
	p.position += stringLength
}

func (p *Outgoing) HeaderAddSize() {
	size := uint16(p.size())
	binary.LittleEndian.PutUint16(p.buffer[p.header-2:], size)
	p.header -= 2
}

func (p *Outgoing) addPadding() {
	size := p.size()
	if size%8 != 0 {
		toAdd := 8 - (size % 8)
		for i := 0; i < toAdd; i++ {
			p.AddUint8(0x33)
		}
	}
}

func (p *Outgoing) XteaEncrypt(xteaKey [4]uint32) error {

	p.HeaderAddSize()
	p.addPadding()

	expandedXteaKey := crypt.ExpandXteaKey(xteaKey)
	crypt.XteaEncrypt(p.buffer[p.header:], expandedXteaKey)

	return nil
}
