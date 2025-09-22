package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Server_InitServer_Success(t *testing.T) {
	myPort = "1234"
	node := &MockNodeAPI{}
	network, err := NewUDPNetwork(":" + myPort)
	assert.NoError(t, err)
	server, err := InitServer(node, network)
	assert.NoError(t, err)
	assert.NotNil(t, server)
	assert.NotNil(t, server.incoming)
	assert.NotNil(t, server.outgoing)
}

/*
func Test_Server_Channel_Operations(t *testing.T) {
	myPort = "4321"
	node := &MockNodeAPI{}
	network, err := NewUDPNetwork(":" + myPort)
	assert.NoError(t, err)
	server, err := InitServer(node, network)
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

func Test_Server_ProcessRequest_STORE(t *testing.T) {
	myPort = "4322"
	node := &MockNodeAPI{}
	network, err := NewUDPNetwork(":" + myPort)
	assert.NoError(t, err)
	server, err := InitServer(node, network)
	assert.NoError(t, err)
	addr := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9998}
	rpc := RPCMessage{Type: "STORE", Payload: Payload{Key: "key", Data: []byte("value"), SourceContact: node.GetSelfContact()}}
	in := IncomingRPC{RPC: rpc, Addr: addr}
	go server.processRequest(in)
	// Check outgoing channel for STORE response
	select {
	case out := <-server.outgoing:
		assert.Equal(t, "STORE", out.RPC.Type)
		assert.Equal(t, "key", out.RPC.Payload.Key)
		assert.Equal(t, node.GetSelfContact(), out.RPC.Payload.Contacts[0])
	case <-time.After(500 * 1e6):
		t.Error("No STORE response received")
	}
}

func Test_Server_ProcessRequest_FIND_VALUE(t *testing.T) {
	myPort = "4323"
	node := &MockNodeAPI{}
	network, err := NewUDPNetwork(":" + myPort)
	assert.NoError(t, err)
	server, err := InitServer(node, network)
	assert.NoError(t, err)
	addr := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9997}
	rpc := RPCMessage{Type: "FIND_VALUE", Payload: Payload{Key: "key", SourceContact: node.GetSelfContact()}}
	in := IncomingRPC{RPC: rpc, Addr: addr}
	go server.processRequest(in)
	// Check outgoing channel for FIND_VALUE response
	select {
	case out := <-server.outgoing:
		assert.Equal(t, "FIND_VALUE", out.RPC.Type)
		// Data is nil in mock
		assert.Nil(t, out.RPC.Payload.Data)
	case <-time.After(500 * 1e6):
		t.Error("No FIND_VALUE response received")
	}
}

func Test_Server_ProcessRequest_Default_Error(t *testing.T) {
	myPort = "4324"
	node := &MockNodeAPI{}
	network, err := NewUDPNetwork(":" + myPort)
	assert.NoError(t, err)
	server, err := InitServer(node, network)
	assert.NoError(t, err)
	addr := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9996}
	rpc := RPCMessage{Type: "UNKNOWN", Payload: Payload{SourceContact: node.GetSelfContact(), TargetContact: node.GetSelfContact()}}
	in := IncomingRPC{RPC: rpc, Addr: addr}
	go server.processRequest(in)
	// Check outgoing channel for ERROR response
	select {
	case out := <-server.outgoing:
		assert.Equal(t, "ERROR", out.RPC.Type)
		assert.Equal(t, node.GetSelfContact(), out.RPC.Payload.TargetContact)
	case <-time.After(500 * 1e6):
		t.Error("No ERROR response received")
	}
}
*/
