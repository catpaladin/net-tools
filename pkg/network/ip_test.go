package network

import (
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockNetworkInterface is a mock implementation of the NetworkInterface.
type MockNetworkInterface struct {
	InterfacesFunc func() ([]net.Interface, error)
	AddrsFunc      func(iface net.Interface) ([]net.Addr, error)
}

func (m MockNetworkInterface) Interfaces() ([]net.Interface, error) {
	return m.InterfacesFunc()
}

func (m MockNetworkInterface) Addrs(iface net.Interface) ([]net.Addr, error) {
	return m.AddrsFunc(iface)
}

func TestGetPrivateIP(t *testing.T) {
	tests := []struct {
		name       string
		interfaces []net.Interface
		addrs      map[string][]net.Addr
		err        error
		expected   string
		expectErr  bool
	}{
		{
			name: "valid IP",
			interfaces: []net.Interface{
				{Name: "eth0", Flags: net.FlagUp},
			},
			addrs: map[string][]net.Addr{
				"eth0": {
					&net.IPNet{IP: net.IPv4(192, 168, 1, 1)},
				},
			},
			expected:  "192.168.1.1",
			expectErr: false,
		},
		{
			name: "loopback IP",
			interfaces: []net.Interface{
				{Name: "lo", Flags: net.FlagUp | net.FlagLoopback},
			},
			addrs: map[string][]net.Addr{
				"lo": {
					&net.IPNet{IP: net.IPv4(127, 0, 0, 1)},
				},
			},
			expected:  "",
			expectErr: true,
		},
		{
			name:       "no interfaces",
			interfaces: []net.Interface{},
			addrs:      map[string][]net.Addr{},
			expected:   "",
			expectErr:  true,
		},
		{
			name: "interface down",
			interfaces: []net.Interface{
				{Name: "eth0", Flags: 0},
			},
			addrs:     map[string][]net.Addr{},
			expected:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNetIf := MockNetworkInterface{
				InterfacesFunc: func() ([]net.Interface, error) {
					return tt.interfaces, tt.err
				},
				AddrsFunc: func(iface net.Interface) ([]net.Addr, error) {
					return tt.addrs[iface.Name], tt.err
				},
			}

			result, err := getPrivateIP(mockNetIf)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// MockHTTPClient is a mock implementation of the HTTPClient interface.
type MockHTTPClient struct {
	GetFunc func(url string) (*http.Response, error)
}

func (m MockHTTPClient) Get(url string) (*http.Response, error) {
	return m.GetFunc(url)
}

// MockResponse is a helper function to create a mock HTTP response.
func MockResponse(body string, statusCode int) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestGetExternalIP(t *testing.T) {
	tests := []struct {
		name      string
		response  *http.Response
		err       error
		expected  string
		expectErr bool
	}{
		{
			name:      "valid IP",
			response:  MockResponse("93.184.216.34\n", http.StatusOK),
			err:       nil,
			expected:  "93.184.216.34",
			expectErr: false,
		},
		{
			name:      "HTTP error",
			response:  nil,
			err:       assert.AnError,
			expected:  "",
			expectErr: true,
		},
		{
			name:      "read body error",
			response:  &http.Response{Body: io.NopCloser(&errorReader{})},
			err:       nil,
			expected:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := MockHTTPClient{
				GetFunc: func(url string) (*http.Response, error) {
					return tt.response, tt.err
				},
			}

			result, err := getExternalIP(mockClient)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// errorReader is a helper type that always returns an error when reading.
type errorReader struct{}

func (e *errorReader) Read(p []byte) (int, error) {
	return 0, assert.AnError
}
