package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_kademlia_InitKademlia_Bootstrap(t *testing.T) {
	k, err := InitKademlia("8000", true, "")
	assert.NoError(t, err)
	assert.NotNil(t, k)
	assert.NotNil(t, k.Node)
	assert.NotNil(t, k.Client)
	assert.NotNil(t, k.Server)
	assert.Equal(t, k.Node.Id.String(), "0000000000000000000000000000000000000000")
}

func Test_kademlia_InitKademlia_Peer(t *testing.T) {
	k, err := InitKademlia("8001", false, "localhost:8000")
	assert.NoError(t, err)
	assert.NotNil(t, k)
	assert.NotNil(t, k.Node)
	assert.NotNil(t, k.Client)
	assert.NotNil(t, k.Server)
	assert.NotEqual(t, k.Node.Id.String(), "0000000000000000000000000000000000000000")
}

func Test_kademlia_GetLocalIP(t *testing.T) {
	ip := GetLocalIP()
	// Should return a non-empty string (may be 127.0.0.1 or actual IP)
	assert.NotEmpty(t, ip)
}
