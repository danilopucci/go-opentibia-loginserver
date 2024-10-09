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

const MAX_INPUT_PACKET_SIZE = 1024

func main() {

	rsaDecrypter, err := crypt.NewRSADecrypter("key.pem")
	if err != nil {
		fmt.Println("Error loading private key:", err)
		os.Exit(1)
	}

	database, err := CreateDatabaseConnection("database_user", "database_password", "127.0.0.1", 3306, "database_name")
	if err != nil {
		fmt.Printf("error while creating database connection: %s\n", err)
	}

	loginParser := NewLoginParser(rsaDecrypter)

	port := 7171
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer listener.Close()

	fmt.Printf("Server is listening on port %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		go handleClient(conn, loginParser, database)
	}

}

func handleClient(conn net.Conn, loginParser *LoginParser, database *sql.DB) {
	defer conn.Close()

	packet := packet.NewIncoming(MAX_INPUT_PACKET_SIZE)

	reqLen, err := conn.Read(packet.PeekBuffer())
	if err != nil {
		fmt.Println("[handleClient] - error reading:", err.Error())
		return
	}
	packet.Resize(reqLen)

	ip, err := utils.GetRemoteIpAddr(conn)
	if err != nil {
		fmt.Printf("[handleClient] - could not get remote IP address: %s\n", err)
		return
	}

	messageSize := packet.GetUint16()
	clientOpcode := Opcode(packet.GetUint8())

	fmt.Printf("[handleClient] - debug - new message (%d bytes); clientOpcode: %d\n", messageSize, clientOpcode)

	switch clientOpcode {
	case Login:
		loginInfo, err := loginParser.ParseLogin(packet)
		if err != nil {
			fmt.Println("[handleClient] - error parsing login info:", err)
		}

		banInfo, err := getIpBanInfo(database, ip)
		if err != nil {
			fmt.Println("[handleClient] - could not fetch ban info:", err)
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
			fmt.Println("[handleClient] - could not fetch account info:", err)
		}

		if utils.Sha1Hash(loginInfo.password) != accountInfo.passwordSHA1 {
			sendClientError(conn, loginInfo.xteaKey, "Account number of password is not correct.")
			return
		}

		accountInfo.characters, err = getCharactersList(database, accountInfo.id)
		if err != nil {
			fmt.Println("[handleClient] - could not fetch character list:", err)
		}

		sendClientMotdAndCharacterList(conn, loginInfo.xteaKey, "Welcome to Test", &accountInfo)
	case Status:
		fmt.Println("status packet is not supported yet")
	}

}

func sendClientError(conn net.Conn, xteaKey [4]uint32, errorData string) {
	packet := packet.NewOutgoing(1024)
	packet.AddUint8(0x0A)
	packet.AddString(errorData)

	sendData(conn, xteaKey, packet)
}

func sendClientMotdAndCharacterList(conn net.Conn, xteaKey [4]uint32, motd string, accountInfo *AccountInfo) {
	packet := packet.NewOutgoing(1024)

	// motd
	if motd != "" {
		packet.AddUint8(0x14)
		packet.AddString(fmt.Sprintf("%d\n%s", 1, motd))
	}

	// character list
	packet.AddUint8(0x64)
	characterListLength := len(accountInfo.characters)
	packet.AddUint8(uint8(characterListLength))

	for i := 0; i < characterListLength; i++ {
		packet.AddString(accountInfo.characters[i])
		packet.AddString("Test")
		packet.AddUint32(2130706433)
		packet.AddUint16(7171)
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
