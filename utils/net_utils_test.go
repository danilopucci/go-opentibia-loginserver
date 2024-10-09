package utils

import (
	"net"
	"testing"
	"time"
)

type mockConn struct {
	remoteAddr string
}

func (m *mockConn) RemoteAddr() net.Addr {
	return &mockAddr{m.remoteAddr}
}

type mockAddr struct {
	addr string
}

func (m *mockAddr) Network() string {
	return "tcp"
}

func (m *mockAddr) String() string {
	return m.addr
}

func (m *mockConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (m *mockConn) Write(b []byte) (n int, err error)  { return 0, nil }
func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestGetRemoteIpAddr(t *testing.T) {
	tests := []struct {
		remoteAddr string
		expectedIP uint32
		expectErr  bool
	}{
		{"192.168.1.1:8080", 3232235777, false}, // 192.168.1.1
		{"10.0.0.1:5000", 167772161, false},     // 10.0.0.1
		{"invalid-ip:1234", 0, true},            // Invalid IP
		{"[::1]:8080", 0, true},                 // IPv6 address
	}

	for _, test := range tests {
		conn := &mockConn{remoteAddr: test.remoteAddr}
		ip, err := GetRemoteIpAddr(conn)
		if test.expectErr {
			if err == nil {
				t.Errorf("expected an error for IP: %s, but got none", test.remoteAddr)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if ip != test.expectedIP {
				t.Errorf("expected %d, got %d for IP %s", test.expectedIP, ip, test.remoteAddr)
			}
		}
	}
}

func TestIpToUint32(t *testing.T) {
	tests := []struct {
		ipStr      string
		expectedIP uint32
		expectErr  bool
	}{
		{"192.168.1.1", 3232235777, false}, // 192.168.1.1
		{"10.0.0.1", 167772161, false},     // 10.0.0.1
		{"invalid-ip", 0, true},            // Invalid IP
		{"::1", 0, true},                   // IPv6 address
	}

	for _, test := range tests {
		ip, err := ipToUint32(test.ipStr)
		if test.expectErr {
			if err == nil {
				t.Errorf("expected an error for IP: %s, but got none", test.ipStr)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if ip != test.expectedIP {
				t.Errorf("expected %d, got %d for IP %s", test.expectedIP, ip, test.ipStr)
			}
		}
	}
}
