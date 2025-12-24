package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	// Simple test to ensure color aliases are callable
	assert.NotEmpty(t, SuccessMsg("test"))
	assert.NotEmpty(t, ErrorMsg("test"))
	assert.NotEmpty(t, DataMsg("test"))
}
