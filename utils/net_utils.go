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

	ip, err := ipToUint32(ipStr)
	if err != nil {
		return 0, fmt.Errorf("could not convert remote IP address (%s): %s", ipStr, err)
	}

	return ip, nil
}

func ipToUint32(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	ip = ip.To4()
	if ip == nil {
		return 0, fmt.Errorf("not an IPv4 address: %s", ipStr)
	}

	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3]), nil
}
