//go:build darwin
// +build darwin

package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatAddress(t *testing.T) {
	tests := []struct {
		ip       string
		expected string
	}{
		{"", "*"},
		{"0.0.0.0", "*"},
		{"::", "[::]"},
		{"127.0.0.1", "127.0.0.1"},
		{"2001:db8::1", "[2001:db8::1]"},
	}

	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			result := formatAddress(tt.ip)
			assert.Equal(t, tt.expected, result)
		})
	}
}
