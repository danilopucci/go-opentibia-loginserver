package utils

import (
	"fmt"
	"net"
)

func GetRemoteIpAddr(conn net.Conn) (uint32, error) {
	clientAddr := conn.RemoteAddr().String()
	ipStr, _, err := net.SplitHostPort(clientAddr)
	if err != nil {
		return 0, fmt.Errorf("could not get remote IP address: %s", err)
	}

	ip, err := IpToUint32(ipStr)
	if err != nil {
		return 0, fmt.Errorf("could not convert remote IP address (%s): %s", ipStr, err)
	}

	return ip, nil
}

func IpToUint32(ipStr string) (uint32, error) {
	if ipStr == "localhost" {
		ipStr = "127.0.0.1"
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	ip = ip.To4()
	if ip == nil {
		return 0, fmt.Errorf("not an IPv4 address: %s", ipStr)
	}

	return uint32(ip[3])<<24 | uint32(ip[2])<<16 | uint32(ip[1])<<8 | uint32(ip[0]), nil
}
