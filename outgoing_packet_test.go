package main

import (
	"encoding/binary"
	"testing"
)

func TestOutgoingPacketInit(t *testing.T) {
	var packet OutgoingPacket
	packet.init(64)

	if len(packet.buffer) != 64 {
		t.Errorf("Expected buffer length of 64, got %d", len(packet.buffer))
	}

	if packet.position != 0 {
		t.Errorf("Expected position to be 0, got %d", packet.position)
	}

	if packet.header != HEADER_OFFSET {
		t.Errorf("Expected header to be %d, got %d", HEADER_OFFSET, packet.header)
	}
}

func TestAddUint8(t *testing.T) {
	var packet OutgoingPacket
	packet.init(15)
	packet.addUint8(0xAB)

	if packet.position != 1 {
		t.Errorf("Expected position to be 1, got %d", packet.position)
	}

	if packet.buffer[HEADER_OFFSET] != 0xAB {
		t.Errorf("Expected buffer to have 0xAB, got %x", packet.buffer[HEADER_OFFSET])
	}
}

func TestAddUint16(t *testing.T) {
	var packet OutgoingPacket
	packet.init(15)
	packet.addUint16(0x1234)

	if packet.position != 2 {
		t.Errorf("Expected position to be 2, got %d", packet.position)
	}

	val := binary.LittleEndian.Uint16(packet.buffer[HEADER_OFFSET:])
	if val != 0x1234 {
		t.Errorf("Expected buffer to have 0x1234, got %x", val)
	}
}

func TestAddUint32(t *testing.T) {
	var packet OutgoingPacket
	packet.init(15)
	packet.addUint32(0xDEADBEEF)

	if packet.position != 4 {
		t.Errorf("Expected position to be 4, got %d", packet.position)
	}

	val := binary.LittleEndian.Uint32(packet.buffer[HEADER_OFFSET:])
	if val != 0xDEADBEEF {
		t.Errorf("Expected buffer to have 0xDEADBEEF, got %x", val)
	}
}

func TestAddString(t *testing.T) {
	var packet OutgoingPacket
	packet.init(64)
	packet.addString("test")

	expectedLen := uint16(4) // Length of "test"
	actualLen := binary.LittleEndian.Uint16(packet.buffer[HEADER_OFFSET:])
	if expectedLen != actualLen {
		t.Errorf("Expected length %d, got %d", expectedLen, actualLen)
	}

	if packet.position != 6 {
		t.Errorf("Expected position to be 6, got %d", packet.position)
	}

	strInBuffer := string(packet.buffer[HEADER_OFFSET+2 : HEADER_OFFSET+6])
	if strInBuffer != "test" {
		t.Errorf("Expected string 'test', got '%s'", strInBuffer)
	}
}

func TestAddPadding(t *testing.T) {
	var packet OutgoingPacket
	packet.init(64)
	packet.addUint8(0xAB)
	packet.addPadding()

	expectedSize := 8 // Should round up to the nearest multiple of 8
	if packet.size() != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, packet.size())
	}

	for i := 1; i < expectedSize; i++ {
		if packet.buffer[HEADER_OFFSET+i] != 0x33 {
			t.Errorf("Expected padding byte 0x33, got %x at position %d", packet.buffer[HEADER_OFFSET+i], i)
		}
	}
}

func TestHeaderAddSize(t *testing.T) {
	var packet OutgoingPacket
	packet.init(64)
	packet.addUint8(0xAB)
	packet.headerAddSize()

	sizeInHeader := binary.LittleEndian.Uint16(packet.buffer[packet.header:])
	if sizeInHeader != 1 {
		t.Errorf("Expected header size 1, got %d", sizeInHeader)
	}

	if packet.header != 8 {
		t.Errorf("Expected header to be 8 after adding size, got %d", packet.header)
	}
}

func TestOutgoingPacketXteaEncrypt(t *testing.T) {
	var packet OutgoingPacket
	packet.init(64)
	packet.addUint32(0xDEADBEEF)

	// Dummy key for XTEA encryption
	xteaKey := [4]uint32{0x1, 0x2, 0x3, 0x4}
	err := packet.xteaEncrypt(xteaKey)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Since actual encryption is not shown, test that the size and padding are handled
	expectedSize := 8 // size should be padded to 8
	if packet.size() != expectedSize {
		t.Errorf("Expected size %d after encryption, got %d", expectedSize, packet.size())
	}
}
