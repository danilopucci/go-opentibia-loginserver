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

const Login uint8 = 0x01
const PACKET_SIZE = 1024

func main() {

	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v", err)
	}

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

		go handleTcpRequest(tcpConnection, loginParser, database, &config)
	}

}

func handleTcpRequest(conn net.Conn, loginParser *LoginParser, database *sql.DB, config *Config) {
	defer conn.Close()

	packet := packet.NewIncoming(PACKET_SIZE)

	reqLen, err := conn.Read(packet.PeekBuffer())
	if err != nil {
		fmt.Printf("[handleClient] - error reading: %s\n", err.Error())
		return
	}
	packet.Resize(reqLen)

	remoteIpAddress, err := utils.GetRemoteIpAddr(conn)
	if err != nil {
		fmt.Printf("[handleClient] - could not get remote IP address: %s\n", err)
		return
	}

	packet.GetUint16() //message size
	clientOpcode := packet.GetUint8()

	if clientOpcode == Login {
		handleLoginRequest(conn, loginParser, database, config, packet, remoteIpAddress)
	} else {
		fmt.Printf("received invalid ClientOpCode (%d) from IP %d\n", clientOpcode, remoteIpAddress)
	}
}

func handleLoginRequest(conn net.Conn, loginParser *LoginParser, database *sql.DB, config *Config, packet *packet.Incoming, remoteIpAddress uint32) {
	loginInfo, err := loginParser.ParseLogin(packet)
	if err != nil {
		fmt.Printf("[handleClient] - error parsing login info: %s\n", err)
		return
	}

	banInfo, err := getIpBanInfo(database, remoteIpAddress)
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
