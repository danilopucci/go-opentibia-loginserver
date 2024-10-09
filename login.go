package main

import (
	"fmt"
	"go-opentibia-loginserver/crypt"
	"go-opentibia-loginserver/packet"
)

type LoginParser struct {
	decrypter crypt.Decrypter
}

type LoginRequest struct {
	clientOs        uint16
	protocolVersion uint16
	datSignature    uint32
	sprSignature    uint32
	picSignature    uint32
	xteaKey         [4]uint32
	accountNumber   uint32
	password        string
}

func NewLoginParser(decrypter crypt.Decrypter) *LoginParser {
	return &LoginParser{decrypter: decrypter}
}

func (loginParser *LoginParser) ParseLogin(packet *packet.Incoming) (LoginRequest, error) {
	var request LoginRequest

	request.clientOs = packet.GetUint16()
	request.protocolVersion = packet.GetUint16()
	request.datSignature = packet.GetUint32()
	request.sprSignature = packet.GetUint32()
	request.picSignature = packet.GetUint32()

	decryptedMsg, err := loginParser.decrypter.DecryptNoPadding(packet.PeekBuffer())
	if err != nil {
		return request, fmt.Errorf("[parseLogin] - error while decrypting packet: %w", err)
	}

	copy(packet.PeekBuffer(), decryptedMsg)

	if packet.GetUint8() != 0 {
		return request, fmt.Errorf("[parseLogin] - error decrypted packet's first byte is not zero")
	}

	request.xteaKey[0] = packet.GetUint32()
	request.xteaKey[1] = packet.GetUint32()
	request.xteaKey[2] = packet.GetUint32()
	request.xteaKey[3] = packet.GetUint32()

	request.accountNumber = packet.GetUint32()
	request.password = packet.GetString()

	return request, nil
}
