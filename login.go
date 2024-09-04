package main

import "fmt"

type LoginHandler struct {
	rsa *RSA
}

func (loginHandler *LoginHandler) init(rsa *RSA) {
	loginHandler.rsa = rsa
}

func (loginHandler *LoginHandler) handleLogin(packet *Packet) {
	clientOs := packet.getUint16()
	protocolVersion := packet.getUint16()
	datSignature := packet.getUint32()
	sprSignature := packet.getUint32()
	picSignature := packet.getUint32()

	fmt.Printf("clientOs: %d; protocolVersion: %d; datSignature: %X; sprSignature: %X, picSignature: %X\n",
		clientOs, protocolVersion, datSignature, sprSignature, picSignature)

	// Decrypt the message
	decryptedMsg, err := loginHandler.rsa.DecryptNoPadding(packet.buffer)
	if err != nil {
		fmt.Println("Error decrypting message:", err)
		return
	}

	packet.buffer = decryptedMsg

	if packet.getUint8() != 0 {
		fmt.Println("Error decrypted is not zero:")
	}

	xteaKey0 := packet.getUint32()
	xteaKey1 := packet.getUint32()
	xteaKey2 := packet.getUint32()
	xteaKey3 := packet.getUint32()
	accountNumber := packet.getUint32()
	password := packet.getString()

	fmt.Printf("xtea0: %d, xtea1: %d, xtea2 %d, xtea3 %d, Account number: %d, password: %s\n", xteaKey0, xteaKey1, xteaKey2, xteaKey3, accountNumber, password)

}
