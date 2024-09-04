package main

import (
	"fmt"
	"net"
	"os"
)

type Opcode uint8

const (
	Login  Opcode = 0x01
	Status Opcode = 0xFF
)

func main() {

	rsaObj := &RSA{}

	// Load the RSA private key from a PEM file
	err := rsaObj.LoadPEM("key.pem")
	if err != nil {
		fmt.Println("Error loading private key:", err)
		os.Exit(1)
	}

	var loginHandler LoginHandler
	loginHandler.init(rsaObj)

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

		go handleClient(conn, &loginHandler)
	}

}

func handleClient(conn net.Conn, loginHandler *LoginHandler) {
	defer conn.Close()

	var packet Packet
	packet.init(1024)

	reqLen, err := conn.Read(packet.buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	packet.resize(reqLen)

	fmt.Printf("Received %d bytes\n", reqLen)

	messageSize := packet.getUint16()
	clientOpcode := Opcode(packet.getUint8())

	fmt.Printf("messageSize: %d; clientOpcode: %d\n", messageSize, clientOpcode)

	switch clientOpcode {
	case Login:
		loginHandler.handleLogin(&packet)
	case Status:
		fmt.Println("status packet is not supported yet")
	}

}
