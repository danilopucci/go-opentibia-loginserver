package protocol

import (
	"fmt"
	"go-opentibia-loginserver/crypt"
	"go-opentibia-loginserver/packet"
)

type LoginParser struct {
	decrypter crypt.Decrypter
}

type LoginRequest struct {
	ClientOs        uint16
	ProtocolVersion uint16
	DatSignature    uint32
	SprSignature    uint32
	PicSignature    uint32
	XteaKey         [4]uint32
	AccountNumber   uint32
	Password        string
}

func NewLoginParser(decrypter crypt.Decrypter) *LoginParser {
	return &LoginParser{decrypter: decrypter}
}

func (loginParser *LoginParser) ParseLogin(packet *packet.Incoming) (LoginRequest, error) {
	var request LoginRequest

	request.ClientOs = packet.GetUint16()
	request.ProtocolVersion = packet.GetUint16()
	request.DatSignature = packet.GetUint32()
	request.SprSignature = packet.GetUint32()
	request.PicSignature = packet.GetUint32()

	decryptedMsg, err := loginParser.decrypter.DecryptNoPadding(packet.PeekBuffer())
	if err != nil {
		return request, fmt.Errorf("[parseLogin] - error while decrypting packet: %w", err)
	}

	copy(packet.PeekBuffer(), decryptedMsg)

	if packet.GetUint8() != 0 {
		return request, fmt.Errorf("[parseLogin] - error decrypted packet's first byte is not zero")
	}

	request.XteaKey[0] = packet.GetUint32()
	request.XteaKey[1] = packet.GetUint32()
	request.XteaKey[2] = packet.GetUint32()
	request.XteaKey[3] = packet.GetUint32()

	request.AccountNumber = packet.GetUint32()
	request.Password = packet.GetString()

	return request, nil
}
