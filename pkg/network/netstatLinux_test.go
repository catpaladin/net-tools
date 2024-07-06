package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockFileSystem is a mock implementation of the FileSystem interface.
type MockFileSystem struct {
	GlobFunc     func(pattern string) ([]string, error)
	ReadlinkFunc func(name string) (string, error)
	ReadFileFunc func(name string) ([]byte, error)
}

func (fs MockFileSystem) Glob(pattern string) ([]string, error) {
	return fs.GlobFunc(pattern)
}

func (fs MockFileSystem) Readlink(name string) (string, error) {
	return fs.ReadlinkFunc(name)
}

func (fs MockFileSystem) ReadFile(name string) ([]byte, error) {
	return fs.ReadFileFunc(name)
}

func TestGetProcessByInode(t *testing.T) {
	tests := []struct {
		name        string
		inode       string
		mockFs      MockFileSystem
		expectedPid int
		expectedCmd string
	}{
		{
			name:  "valid inode",
			inode: "12345",
			mockFs: MockFileSystem{
				GlobFunc: func(pattern string) ([]string, error) {
					return []string{
						"/proc/123/fd/4",
					}, nil
				},
				ReadlinkFunc: func(name string) (string, error) {
					if name == "/proc/123/fd/4" {
						return "socket:[12345]", nil
					}
					return "", nil
				},
				ReadFileFunc: func(name string) ([]byte, error) {
					if name == "/proc/123/cmdline" {
						return []byte("bash\x00-c\x00ls"), nil
					}
					return nil, nil
				},
			},
			expectedPid: 123,
			expectedCmd: "bash -c ls",
		},
		{
			name:  "inode not found",
			inode: "12345",
			mockFs: MockFileSystem{
				GlobFunc: func(pattern string) ([]string, error) {
					return []string{}, nil
				},
				ReadlinkFunc: func(name string) (string, error) {
					return "", nil
				},
				ReadFileFunc: func(name string) ([]byte, error) {
					return nil, nil
				},
			},
			expectedPid: 0,
			expectedCmd: "",
		},
		{
			name:  "invalid inode format",
			inode: "12345",
			mockFs: MockFileSystem{
				GlobFunc: func(pattern string) ([]string, error) {
					return []string{
						"/proc/invalid/fd/4",
					}, nil
				},
				ReadlinkFunc: func(name string) (string, error) {
					return "", nil
				},
				ReadFileFunc: func(name string) ([]byte, error) {
					return nil, nil
				},
			},
			expectedPid: 0,
			expectedCmd: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pid, cmd := getProcessByInode(tt.mockFs, tt.inode)
			assert.Equal(t, tt.expectedPid, pid)
			assert.Equal(t, tt.expectedCmd, cmd)
		})
	}
}

func TestParseLinuxTCPLine(t *testing.T) {
	tests := []struct {
		name               string
		line               string
		expectedLocalAddr  string
		expectedLocalPort  string
		expectedRemoteAddr string
		expectedRemotePort string
		expectedStatus     string
		expectedInode      string
	}{
		{
			name:               "valid line",
			line:               "  5: 0100007F:0035 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000 0 123456 1 ffff8881b86ad400 99 0 0 10 0",
			expectedLocalAddr:  "127.0.0.1",
			expectedLocalPort:  "53",
			expectedRemoteAddr: "0.0.0.0",
			expectedRemotePort: "0",
			expectedStatus:     "0A",
			expectedInode:      "123456",
		},
		{
			name:               "invalid line",
			line:               "invalid line",
			expectedLocalAddr:  "",
			expectedLocalPort:  "",
			expectedRemoteAddr: "",
			expectedRemotePort: "",
			expectedStatus:     "",
			expectedInode:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localAddr, localPort, remoteAddr, remotePort, status, inode := parseLinuxTCPLine(tt.line)
			assert.Equal(t, tt.expectedLocalAddr, localAddr)
			assert.Equal(t, tt.expectedLocalPort, localPort)
			assert.Equal(t, tt.expectedRemoteAddr, remoteAddr)
			assert.Equal(t, tt.expectedRemotePort, remotePort)
			assert.Equal(t, tt.expectedStatus, status)
			assert.Equal(t, tt.expectedInode, inode)
		})
	}
}

func TestParseHexIPPort(t *testing.T) {
	tests := []struct {
		name         string
		hex          string
		expectedIP   string
		expectedPort string
	}{
		{
			name:         "valid IPv4",
			hex:          "0100007F:0035",
			expectedIP:   "127.0.0.1",
			expectedPort: "53",
		},
		{
			name:         "invalid hex",
			hex:          "invalid",
			expectedIP:   "",
			expectedPort: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, port := parseHexIPPort(tt.hex)
			assert.Equal(t, tt.expectedIP, ip)
			assert.Equal(t, tt.expectedPort, port)
		})
	}
}

func TestParseHexIP(t *testing.T) {
	tests := []struct {
		name     string
		hex      string
		expected string
	}{
		{
			name:     "valid IPv4",
			hex:      "0100007F",
			expected: "127.0.0.1",
		},
		{
			name:     "valid IPv6",
			hex:      "00000000000000000000000000000001",
			expected: "::1",
		},
		{
			name:     "invalid hex",
			hex:      "invalid",
			expected: "0.0.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHexIP(tt.hex)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseHexIPv4(t *testing.T) {
	tests := []struct {
		name     string
		hex      string
		expected string
	}{
		{
			name:     "valid IPv4",
			hex:      "0100007F",
			expected: "127.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHexIPv4(tt.hex)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseHexIPv6(t *testing.T) {
	tests := []struct {
		name     string
		hex      string
		expected string
	}{
		{
			name:     "valid IPv6",
			hex:      "00000000000000000000000000000001",
			expected: "::1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHexIPv6(tt.hex)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseHexPort(t *testing.T) {
	tests := []struct {
		name     string
		hex      string
		expected string
	}{
		{
			name:     "valid port",
			hex:      "0035",
			expected: "53",
		},
		{
			name:     "invalid hex",
			hex:      "invalid",
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHexPort(tt.hex)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseHexToInt(t *testing.T) {
	tests := []struct {
		name     string
		hex      string
		expected int
	}{
		{
			name:     "valid hex",
			hex:      "0035",
			expected: 53,
		},
		{
			name:     "invalid hex",
			hex:      "invalid",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHexToInt(tt.hex)
			assert.Equal(t, tt.expected, result)
		})
	}
}
