package main

import (
	"bytes"
	"encoding/binary"
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

	message := make([]byte, 1024)
	reqLen, err := conn.Read(message)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	message = message[:reqLen]

	fmt.Printf("Received %d bytes\n", reqLen)

	messageSize := getUint16(&message)
	clientOpcode := getUint8(&message)
	clientOs := getUint16(&message)
	protocolVersion := getUint16(&message)
	datSignature := getUint32(&message)
	sprSignature := getUint32(&message)
	picSignature := getUint32(&message)

	fmt.Printf("messageSize: %d, clientOpcode: %d; clientOs: %d; protocolVersion: %d; datSignature: %X; sprSignature: %X, picSignature: %X\n",
		messageSize, clientOpcode, clientOs, protocolVersion, datSignature, sprSignature, picSignature)

	fmt.Printf("Encrypted message length: %d\n", len(message))
	// Decrypt the message
	decryptedMsg, err := rsaObj.DecryptNoPadding(message)
	if err != nil {
		fmt.Println("Error decrypting message:", err)
		return
	}

	if getUint8(&decryptedMsg) != 0 {
		fmt.Println("Error decrypted is not zero:")
	}

	fmt.Printf("Encrypted message length: %d\n", len(decryptedMsg))

	xteaKey0 := getUint32(&decryptedMsg)
	xteaKey1 := getUint32(&decryptedMsg)
	xteaKey2 := getUint32(&decryptedMsg)
	xteaKey3 := getUint32(&decryptedMsg)
	accountNumber := getUint32(&decryptedMsg)
	password := getString(&decryptedMsg)

	fmt.Printf("xtea0: %d, xtea1: %d, xtea2 %d, xtea3 %d, Account number: %d, password: %s\n", xteaKey0, xteaKey1, xteaKey2, xteaKey3, accountNumber, password)

	fmt.Printf("Encrypted message length: %d; data: %x\n", len(decryptedMsg), decryptedMsg)
}

func skipBytes(msg *[]byte, n int) {
	*msg = (*msg)[n:]
}

// getUint8 reads a uint8 from the message.
func getUint8(msg *[]byte) uint8 {
	var result uint8
	buf := bytes.NewReader(*msg)
	binary.Read(buf, binary.LittleEndian, &result)
	*msg = (*msg)[1:] // Skip the bytes we've just read
	return result
}

// getUint16 reads a uint16 from the message.
func getUint16(msg *[]byte) uint16 {
	var result uint16
	buf := bytes.NewReader(*msg)
	binary.Read(buf, binary.LittleEndian, &result)
	*msg = (*msg)[2:] // Skip the bytes we've just read
	return result
}

// getUint32 reads a uint32 from the message.
func getUint32(msg *[]byte) uint32 {
	var result uint32
	buf := bytes.NewReader(*msg)
	binary.Read(buf, binary.LittleEndian, &result)
	*msg = (*msg)[4:] // Skip the bytes we've just read
	return result
}

// getString reads a string from the message.
func getString(msg *[]byte) string {
	var result string
	stringLength := getUint16(msg)

	result = string((*msg)[:stringLength])
	*msg = (*msg)[stringLength:] // Skip the bytes we've just read

	return result
}
