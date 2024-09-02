package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	rsaObj := &RSA{}

	// Load the RSA private key from a PEM file
	err := rsaObj.LoadPEM("key.pem")
	if err != nil {
		fmt.Println("Error loading private key:", err)
		os.Exit(1)
	}

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

		go handleClient(conn, rsaObj)
	}

}

func handleClient(conn net.Conn, rsaObj *RSA) {
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
	clientOpcode := packet.getUint8()
	clientOs := packet.getUint16()
	protocolVersion := packet.getUint16()
	datSignature := packet.getUint32()
	sprSignature := packet.getUint32()
	picSignature := packet.getUint32()

	fmt.Printf("messageSize: %d, clientOpcode: %d; clientOs: %d; protocolVersion: %d; datSignature: %X; sprSignature: %X, picSignature: %X\n",
		messageSize, clientOpcode, clientOs, protocolVersion, datSignature, sprSignature, picSignature)

	// Decrypt the message
	decryptedMsg, err := rsaObj.DecryptNoPadding(packet.buffer)
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
