package protocol

import (
	"fmt"
	"go-opentibia-loginserver/config"
	"go-opentibia-loginserver/models"
	"go-opentibia-loginserver/packet"
	"go-opentibia-loginserver/utils"
	"net"
)

const PACKET_SIZE = 1024

func SendClientError(conn net.Conn, xteaKey [4]uint32, errorData string) {
	packet := packet.NewOutgoing(PACKET_SIZE)
	packet.AddUint8(0x0A)
	packet.AddString(errorData)

	SendData(conn, xteaKey, packet)
}

func SendClientMotdAndCharacterList(conn net.Conn, xteaKey [4]uint32, motd string, accountInfo *models.AccountInfo, cfg *config.Config) {
	packet := packet.NewOutgoing(PACKET_SIZE)

	// motd
	if motd != "" {
		packet.AddUint8(0x14)
		packet.AddString(fmt.Sprintf("%d\n%s", 1, motd))
	}

	// character list
	packet.AddUint8(0x64)
	characterListLength := len(accountInfo.Characters)
	packet.AddUint8(uint8(characterListLength))

	//there is no support for multiworld yet, so get the default world
	world := config.GetDefaultWorld(cfg)

	for i := 0; i < characterListLength; i++ {
		packet.AddString(accountInfo.Characters[i])
		packet.AddString(world.Name)
		packet.AddUint32(world.HostIP)
		packet.AddUint16(world.Port)
	}

	premiumDays := utils.CalculateRemainingDays(accountInfo.PremiumEndsAt)
	if premiumDays < 0 {
		packet.AddUint16(0)
	} else {
		packet.AddUint16(uint16(premiumDays))
	}

	SendData(conn, xteaKey, packet)
}

func SendData(conn net.Conn, xteaKey [4]uint32, packet *packet.Outgoing) error {
	packet.XteaEncrypt(xteaKey)
	packet.HeaderAddSize()

	dataToSend := packet.Get()

	_, err := conn.Write(dataToSend)
	if err != nil {
		return fmt.Errorf("failed to send data: %v", err)
	}

	return nil
}
