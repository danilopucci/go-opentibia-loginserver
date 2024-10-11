package main

import (
	"database/sql"
	"fmt"
	"go-opentibia-loginserver/config"
	"go-opentibia-loginserver/crypt"
	"go-opentibia-loginserver/database"
	"go-opentibia-loginserver/packet"
	"go-opentibia-loginserver/protocol"
	"go-opentibia-loginserver/utils"
	"net"
	"os"
)

const Login uint8 = 0x01
const PACKET_SIZE = 1024

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v", err)
	}

	rsaDecrypter, err := crypt.NewRSADecrypter(config.RSAKeyFile)
	if err != nil {
		fmt.Println("Error loading private key:", err)
		os.Exit(1)
	}

	db, err := database.CreateDatabaseConnection(config.Database.User, config.Database.Password, config.Database.HostName, config.Database.Port, config.Database.Name)
	if err != nil {
		fmt.Printf("error while creating database connection: %s\n", err)
	}

	loginParser := protocol.NewLoginParser(rsaDecrypter)

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

		go handleTcpRequest(tcpConnection, loginParser, db, &config)
	}

}

func handleTcpRequest(conn net.Conn, loginParser *protocol.LoginParser, db *sql.DB, cfg *config.Config) {
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
		handleLoginRequest(conn, loginParser, db, cfg, packet, remoteIpAddress)
	} else {
		fmt.Printf("received invalid ClientOpCode (%d) from IP %d\n", clientOpcode, remoteIpAddress)
	}
}

func handleLoginRequest(conn net.Conn, loginParser *protocol.LoginParser, db *sql.DB, cfg *config.Config, packet *packet.Incoming, remoteIpAddress uint32) {
	loginInfo, err := loginParser.ParseLogin(packet)
	if err != nil {
		fmt.Printf("[handleClient] - error parsing login info: %s\n", err)
		return
	}

	databaseQuery := database.GetDatabaseQuery(cfg.QueryVersion)

	banInfo, err := databaseQuery.GetIpBanInfo(db, remoteIpAddress)
	if err != nil {
		fmt.Printf("[handleClient] - could not fetch ban info: %s\n", err)
		return
	}

	if banInfo.IsBanned {
		banExpiresDateTime := utils.FormatDateTimeUTC(banInfo.ExpiresAt)
		protocol.SendClientError(conn, loginInfo.XteaKey, fmt.Sprintf("Your IP has been banned until %s.\n\nReason specified:\n%s", banExpiresDateTime, banInfo.Reason))
		return
	}

	if loginInfo.AccountNumber == 0 {
		protocol.SendClientError(conn, loginInfo.XteaKey, "Invalid account number.")
		return
	}

	if loginInfo.Password == "" {
		protocol.SendClientError(conn, loginInfo.XteaKey, "Invalid password.")
		return
	}

	accountInfo, err := databaseQuery.GetAccountInfo(db, loginInfo.AccountNumber)
	if err != nil {
		fmt.Printf("[handleClient] - could not fetch account info: %s\n", err)
		return
	}

	if utils.Sha1Hash(loginInfo.Password) != accountInfo.PasswordSHA1 {
		protocol.SendClientError(conn, loginInfo.XteaKey, "Account number of password is not correct.")
		return
	}

	accountInfo.Characters, err = databaseQuery.GetCharactersList(db, accountInfo.Id)
	if err != nil {
		fmt.Printf("[handleClient] - could not fetch character list: %s\n", err)
		return
	}

	protocol.SendClientMotdAndCharacterList(conn, loginInfo.XteaKey, cfg.Motd, &accountInfo, cfg)
}
