package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommands(t *testing.T) {
	t.Run("dns", func(t *testing.T) {
		cmd := DNSCommand()
		assert.Equal(t, "dns", cmd.Use)
		assert.NotNil(t, cmd.Run)
	})

	t.Run("ip", func(t *testing.T) {
		cmd := IPCmd()
		assert.Equal(t, "ip", cmd.Use)
		assert.NotNil(t, cmd.Run)
	})

	t.Run("port", func(t *testing.T) {
		cmd := PortCmd()
		assert.Equal(t, "port", cmd.Use)
		assert.NotNil(t, cmd.Run)
	})

	t.Run("processes", func(t *testing.T) {
		cmd := ProcessesCmd()
		assert.Equal(t, "processes", cmd.Use)
		assert.NotNil(t, cmd.Run)
	})

	t.Run("tui", func(t *testing.T) {
		cmd := TUICmd()
		assert.Equal(t, "tui", cmd.Use)
		assert.NotNil(t, cmd.Run)
	})
}
