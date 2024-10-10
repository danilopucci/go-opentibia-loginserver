package main

import (
	"database/sql"
	"fmt"
	"go-opentibia-loginserver/crypt"
	"go-opentibia-loginserver/packet"
	"go-opentibia-loginserver/utils"
	"net"
	"os"
)

type Opcode uint8

const (
	Login  Opcode = 0x01
	Status Opcode = 0xFF
)

const PACKET_SIZE = 1024

func main() {

	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v", err)
	}

	for i := range config.GameServer.Worlds {
		ipAddress, err := utils.IpToUint32(config.GameServer.Worlds[i].HostName)
		if err == nil {
			config.GameServer.Worlds[i].HostIP = ipAddress
		} else {
			fmt.Println("could not convert world %s host %s to number ip address: %s", config.GameServer.Worlds[i].Name, config.GameServer.Worlds[i].HostName, err)
		}
	}

	fmt.Println("%s", config.Database.HostName)
	fmt.Println("DATABASE_HOST:", os.Getenv("DATABASE_HOST"))

	rsaDecrypter, err := crypt.NewRSADecrypter(config.RSAKeyFile)
	if err != nil {
		fmt.Println("Error loading private key:", err)
		os.Exit(1)
	}

	database, err := CreateDatabaseConnection(config.Database.User, config.Database.Password, config.Database.HostName, config.Database.Port, config.Database.Name)
	if err != nil {
		fmt.Printf("error while creating database connection: %s\n", err)
	}

	loginParser := NewLoginParser(rsaDecrypter)

	tcpListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.LoginServer.HostName, config.LoginServer.Port))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer tcpListener.Close()

	for {
		tcpConnection, err := tcpListener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go handleLoginRequest(tcpConnection, loginParser, database, &config)
	}

}

func handleLoginRequest(conn net.Conn, loginParser *LoginParser, database *sql.DB, config *Config) {
	defer conn.Close()

	packet := packet.NewIncoming(PACKET_SIZE)

	reqLen, err := conn.Read(packet.PeekBuffer())
	if err != nil {
		fmt.Printf("[handleClient] - error reading: %s\n", err.Error())
		return
	}
	packet.Resize(reqLen)

	ip, err := utils.GetRemoteIpAddr(conn)
	if err != nil {
		fmt.Printf("[handleClient] - could not get remote IP address: %s\n", err)
		return
	}

	packet.GetUint16()
	clientOpcode := Opcode(packet.GetUint8())

	if clientOpcode != Login {
		fmt.Printf("received invalid ClientOpCode (%d) from IP %d\n", clientOpcode, ip)
		return
	}

	loginInfo, err := loginParser.ParseLogin(packet)
	if err != nil {
		fmt.Printf("[handleClient] - error parsing login info: %s\n", err)
		return
	}

	banInfo, err := getIpBanInfo(database, ip)
	if err != nil {
		fmt.Printf("[handleClient] - could not fetch ban info: %s\n", err)
		return
	}

	if banInfo.isBanned {
		banExpiresDateTime := utils.FormatDateTimeUTC(banInfo.expiresAt)
		sendClientError(conn, loginInfo.xteaKey, fmt.Sprintf("Your IP has been banned until %s.\n\nReason specified:\n%s", banExpiresDateTime, banInfo.reason))
		return
	}

	if loginInfo.accountNumber == 0 {
		sendClientError(conn, loginInfo.xteaKey, "Invalid account number.")
		return
	}

	if loginInfo.password == "" {
		sendClientError(conn, loginInfo.xteaKey, "Invalid password.")
		return
	}

	accountInfo, err := getAccountInfo(database, loginInfo.accountNumber)
	if err != nil {
		fmt.Printf("[handleClient] - could not fetch account info: %s\n", err)
		return
	}

	if utils.Sha1Hash(loginInfo.password) != accountInfo.passwordSHA1 {
		sendClientError(conn, loginInfo.xteaKey, "Account number of password is not correct.")
		return
	}

	accountInfo.characters, err = getCharactersList(database, accountInfo.id)
	if err != nil {
		fmt.Printf("[handleClient] - could not fetch character list: %s\n", err)
		return
	}

	sendClientMotdAndCharacterList(conn, loginInfo.xteaKey, config.Motd, &accountInfo, config)
}

func sendClientError(conn net.Conn, xteaKey [4]uint32, errorData string) {
	packet := packet.NewOutgoing(PACKET_SIZE)
	packet.AddUint8(0x0A)
	packet.AddString(errorData)

	sendData(conn, xteaKey, packet)
}

func sendClientMotdAndCharacterList(conn net.Conn, xteaKey [4]uint32, motd string, accountInfo *AccountInfo, config *Config) {
	packet := packet.NewOutgoing(PACKET_SIZE)

	// motd
	if motd != "" {
		packet.AddUint8(0x14)
		packet.AddString(fmt.Sprintf("%d\n%s", 1, motd))
	}

	// character list
	packet.AddUint8(0x64)
	characterListLength := len(accountInfo.characters)
	packet.AddUint8(uint8(characterListLength))

	//there is no support for multiworld yet, so get the default world
	world := GetDefaultWorld(config)

	for i := 0; i < characterListLength; i++ {
		packet.AddString(accountInfo.characters[i])
		packet.AddString(world.Name)
		packet.AddUint32(world.HostIP)
		packet.AddUint16(world.Port)
	}
	packet.AddUint16(20)

	sendData(conn, xteaKey, packet)
}

func sendData(conn net.Conn, xteaKey [4]uint32, packet *packet.Outgoing) error {
	packet.XteaEncrypt(xteaKey)
	packet.HeaderAddSize()

	dataToSend := packet.Get()

	_, err := conn.Write(dataToSend)
	if err != nil {
		return fmt.Errorf("failed to send data: %v", err)
	}

	return nil
}
