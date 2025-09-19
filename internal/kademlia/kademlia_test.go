package kademlia

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_kademlia_InitKademlia_Bootstrap(t *testing.T) {
	k, err := InitKademlia("8000", true, "", WithSkipBootstrapPing(true))
	assert.NoError(t, err)
	assert.NotNil(t, k)
	assert.NotNil(t, k.Node)
	assert.NotNil(t, k.Client)
	assert.NotNil(t, k.Server)
	assert.Equal(t, k.Node.Id.String(), "0000000000000000000000000000000000000000")
}

func Test_kademlia_InitKademlia_Peer(t *testing.T) {
	k, err := InitKademlia("8001", false, "localhost:8000", WithSkipBootstrapPing(true))
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

func Test_kademlia_PingBetweenNodes(t *testing.T) {
	// Use two different ports
	portA := "9100"
	portB := "9101"

	// Start node A (bootstrap)
	nodeA, errA := InitKademlia(portA, true, "", WithSkipBootstrapPing(true))
	assert.NoError(t, errA)
	assert.NotNil(t, nodeA)
	nodeA.Server.RunServer()

	// Start node B (peer)
	nodeB, errB := InitKademlia(portB, false, GetLocalIP()+":"+portA, WithSkipBootstrapPing(true))
	assert.NoError(t, errB)
	assert.NotNil(t, nodeB)
	nodeB.Server.RunServer()

	// Let servers start
	time.Sleep(500 * time.Millisecond)

	// Node B pings node A
	target := nodeA.Node.GetSelfContact()
	resp, err := nodeB.Client.SendPingMessage(target)
	assert.NoError(t, err)
	assert.Equal(t, "PONG", resp.Type)

	// Node A pings node B
	targetB := nodeB.Node.GetSelfContact()
	resp2, err2 := nodeA.Client.SendPingMessage(targetB)
	assert.NoError(t, err2)
	assert.Equal(t, "PONG", resp2.Type)
}
