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

type dummyAddr struct{}

func (d dummyAddr) Network() string { return "tcp" }
func (d dummyAddr) String() string  { return "1.2.3.4" }

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
			name: "loopback IP on non-loopback interface",
			interfaces: []net.Interface{
				{Name: "eth0", Flags: net.FlagUp},
			},
			addrs: map[string][]net.Addr{
				"eth0": {
					&net.IPNet{IP: net.IPv4(127, 0, 0, 1)},
				},
			},
			expected:  "",
			expectErr: true,
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
		{
			name: "interfaces error",
			interfaces: []net.Interface{
				{Name: "eth0", Flags: net.FlagUp},
			},
			addrs: map[string][]net.Addr{
				"eth0": {
					&net.IPNet{IP: net.IPv4(192, 168, 1, 1)},
				},
			},
			err:       assert.AnError,
			expected:  "",
			expectErr: true,
		},
		{
			name: "addrs error",
			interfaces: []net.Interface{
				{Name: "eth0", Flags: net.FlagUp},
			},
			addrs:     map[string][]net.Addr{},
			err:       assert.AnError,
			expected:  "",
			expectErr: true,
		},
		{
			name: "IPAddr type",
			interfaces: []net.Interface{
				{Name: "eth0", Flags: net.FlagUp},
			},
			addrs: map[string][]net.Addr{
				"eth0": {
					&net.IPAddr{IP: net.IPv4(10, 0, 0, 1)},
				},
			},
			expected:  "10.0.0.1",
			expectErr: false,
		},
		{
			name: "IPv6 ignored",
			interfaces: []net.Interface{
				{Name: "eth1", Flags: net.FlagUp},
			},
			addrs: map[string][]net.Addr{
				"eth1": {
					&net.IPNet{IP: net.ParseIP("fe80::1")},
				},
			},
			expected:  "",
			expectErr: true,
		},
		{
			name: "unknown addr type",
			interfaces: []net.Interface{
				{Name: "eth0", Flags: net.FlagUp},
			},
			addrs: map[string][]net.Addr{
				"eth0": {
					dummyAddr{},
				},
			},
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

func TestGetIPWithDependencies(t *testing.T) {
	mockNetIf := MockNetworkInterface{
		InterfacesFunc: func() ([]net.Interface, error) {
			return []net.Interface{{Name: "eth0", Flags: net.FlagUp}}, nil
		},
		AddrsFunc: func(iface net.Interface) ([]net.Addr, error) {
			return []net.Addr{&net.IPNet{IP: net.IPv4(192, 168, 1, 1)}}, nil
		},
	}
	mockClient := MockHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			return MockResponse("1.2.3.4", 200), nil
		},
	}

	t.Run("private", func(t *testing.T) {
		res, err := getIPWithDependencies("private", mockNetIf, mockClient)
		assert.NoError(t, err)
		assert.Equal(t, "192.168.1.1", res)
	})

	t.Run("public", func(t *testing.T) {
		res, err := getIPWithDependencies("public", mockNetIf, mockClient)
		assert.NoError(t, err)
		assert.Equal(t, "1.2.3.4", res)
	})

	t.Run("invalid", func(t *testing.T) {
		_, err := getIPWithDependencies("invalid", mockNetIf, mockClient)
		assert.Error(t, err)
	})

	t.Run("private error", func(t *testing.T) {
		errNetIf := MockNetworkInterface{
			InterfacesFunc: func() ([]net.Interface, error) {
				return nil, assert.AnError
			},
		}
		_, err := getIPWithDependencies("private", errNetIf, mockClient)
		assert.Error(t, err)
	})

	t.Run("public error", func(t *testing.T) {
		errClient := MockHTTPClient{
			GetFunc: func(url string) (*http.Response, error) {
				return nil, assert.AnError
			},
		}
		_, err := getIPWithDependencies("public", mockNetIf, errClient)
		assert.Error(t, err)
	})
}
