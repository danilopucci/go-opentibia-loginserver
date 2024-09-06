package main

import (
	"fmt"
	"net"
)

type LoginHandler struct {
	rsa *RSA
}

func (loginHandler *LoginHandler) init(rsa *RSA) {
	loginHandler.rsa = rsa
}

func (loginHandler *LoginHandler) handleLogin(packet *IncomingPacket, conn net.Conn) {
	clientOs := packet.getUint16()
	protocolVersion := packet.getUint16()
	datSignature := packet.getUint32()
	sprSignature := packet.getUint32()
	picSignature := packet.getUint32()

	fmt.Printf("clientOs: %d; protocolVersion: %d; datSignature: %X; sprSignature: %X, picSignature: %X\n",
		clientOs, protocolVersion, datSignature, sprSignature, picSignature)

	// Decrypt the message
	decryptedMsg, err := loginHandler.rsa.DecryptNoPadding(packet.peekBuffer())
	if err != nil {
		fmt.Println("Error decrypting message:", err)
		return
	}

	copy(packet.buffer[packet.position:], decryptedMsg)

	if packet.getUint8() != 0 {
		fmt.Println("Error decrypted is not zero:")
	}

	var xteaKey [4]uint32
	xteaKey[0] = packet.getUint32()
	xteaKey[1] = packet.getUint32()
	xteaKey[2] = packet.getUint32()
	xteaKey[3] = packet.getUint32()

	accountNumber := packet.getUint32()
	password := packet.getString()

	fmt.Printf("xtea0: %d, xtea1: %d, xtea2 %d, xtea3 %d, Account number: %d, password: %s\n", xteaKey[0], xteaKey[1], xteaKey[2], xteaKey[3], accountNumber, password)

	loginHandler.sendError(conn, xteaKey, "Teste Alala O")
}

func (loginHandler *LoginHandler) sendError(conn net.Conn, xteaKey [4]uint32, errorData string) {
	var packet OutgoingPacket
	packet.init(1024)
	packet.addUint8(0x0A)
	packet.addString(errorData)

	loginHandler.send(conn, xteaKey, &packet)
}

func (loginHandler *LoginHandler) send(conn net.Conn, xteaKey [4]uint32, packet *OutgoingPacket) error {

	//fmt.Printf("raw data before encrypt: header: %d, position: %d, data: %x\n", packet.header, packet.position, packet.buffer)
	packet.xteaEncrypt(xteaKey)
	//fmt.Printf("raw data after encrypt: header: %d, position: %d, data: %x\n", packet.header, packet.position, packet.buffer)
	packet.headerAddSize()
	//fmt.Printf("raw data after header size: header: %d, position: %d, data: %x\n", packet.header, packet.position, packet.buffer)

	dataToSend := packet.get()
	fmt.Printf("data output: %d\n", dataToSend)

	n, err := conn.Write(dataToSend)
	if err != nil {
		return fmt.Errorf("failed to send data: %v", err)
	}

	fmt.Printf("Sent %d bytes\n", n)
	return nil
}
