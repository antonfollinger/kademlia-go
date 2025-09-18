package kademlia

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Server_InitServer_Success(t *testing.T) {
	myPort = "1234"
	node := &MockNodeAPI{}
	server, err := InitServer(node)
	assert.NoError(t, err)
	assert.NotNil(t, server)
	assert.NotNil(t, server.conn)
	assert.NotNil(t, server.incoming)
	assert.NotNil(t, server.outgoing)
}

func Test_Server_Channel_Operations(t *testing.T) {
	myPort = "4321"
	node := &MockNodeAPI{}
	server, err := InitServer(node)
	assert.NoError(t, err)
	// Test sending and receiving on incoming channel
	rpc := RPCMessage{Type: "PING"}
	addr := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9999}
	in := IncomingRPC{RPC: rpc, Addr: addr}
	select {
	case server.incoming <- in:
		// success
	default:
		t.Error("Failed to send to incoming channel")
	}
	// Test sending and receiving on outgoing channel
	out := OutgoingRPC{RPC: rpc, Addr: addr}
	select {
	case server.outgoing <- out:
		// success
	default:
		t.Error("Failed to send to outgoing channel")
	}
}
