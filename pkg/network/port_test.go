package network

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockDialer is a mock implementation of the Dialer interface.
type MockDialer struct {
	DialFunc func(network, address string) (net.Conn, error)
}

func (m MockDialer) Dial(network, address string) (net.Conn, error) {
	return m.DialFunc(network, address)
}

// MockConn is a mock implementation of the net.Conn interface.
type MockConn struct{}

func (m MockConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (m MockConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (m MockConn) Close() error                       { return nil }
func (m MockConn) LocalAddr() net.Addr                { return nil }
func (m MockConn) RemoteAddr() net.Addr               { return nil }
func (m MockConn) SetDeadline(t time.Time) error      { return nil }
func (m MockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m MockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestNetcat(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		port      string
		dialFunc  func(network, address string) (net.Conn, error)
		expectErr bool
	}{
		{
			name: "successful connection",
			host: "localhost",
			port: "8080",
			dialFunc: func(network, address string) (net.Conn, error) {
				assert.Equal(t, "tcp", network)
				assert.Equal(t, "localhost:8080", address)
				return MockConn{}, nil
			},
			expectErr: false,
		},
		{
			name: "connection error",
			host: "localhost",
			port: "8080",
			dialFunc: func(network, address string) (net.Conn, error) {
				return nil, errors.New("connection refused")
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDialer := MockDialer{
				DialFunc: tt.dialFunc,
			}

			err := netcatDialer(mockDialer, tt.host, tt.port)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
