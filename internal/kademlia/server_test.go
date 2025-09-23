package kademlia

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Server_InitServer_Success(t *testing.T) {
	port := "1234"
	node := &MockNodeAPI{Port: port}
	registry := NewMockRegistry()
	network := NewMockNetwork("127.0.0.1:"+port, registry)
	server, err := InitServer(node, network)
	assert.NoError(t, err)
	assert.NotNil(t, server)
	assert.NotNil(t, server.incoming)
	assert.NotNil(t, server.outgoing)
}

func Test_Server_Channel_Operations(t *testing.T) {
	port := "4321"
	node := &MockNodeAPI{Port: port}
	registry := NewMockRegistry()
	network := NewMockNetwork("127.0.0.1:"+port, registry)
	server, err := InitServer(node, network)
	assert.NoError(t, err)
	// Test sending and receiving on incoming channel
	rpc := RPCMessage{Type: "PING"}
	addr := "127.0.0.1:9999"
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
	port := "4322"
	node := &MockNodeAPI{Port: port}
	registry := NewMockRegistry()
	network := NewMockNetwork("127.0.0.1:"+port, registry)
	server, err := InitServer(node, network)
	assert.NoError(t, err)
	addr := "127.0.0.1:9998"
	registry.Register(addr)
	rpc := NewRPCMessage("STORE", Payload{Key: "key", Data: []byte("value"), SourceContact: node.GetSelfContact()}, true)
	in := IncomingRPC{RPC: *rpc, Addr: addr}
	server.incoming <- in
	time.Sleep(500 * time.Millisecond)
	ch, ok := registry.Get(addr)
	assert.True(t, ok)
	select {
	case pkt := <-ch:
		var outRPC RPCMessage
		err := json.Unmarshal(pkt.data, &outRPC)
		assert.NoError(t, err)
		assert.Equal(t, "STORE", outRPC.Type)
		assert.Equal(t, "key", outRPC.Payload.Key)
		assert.Equal(t, node.GetSelfContact(), outRPC.Payload.Contacts[0])
	case <-time.After(1 * time.Second):
		t.Error("No STORE response received")
	}
}

func Test_Server_ProcessRequest_FIND_VALUE(t *testing.T) {
	port := "4323"
	node := &MockNodeAPI{Port: port}
	registry := NewMockRegistry()
	network := NewMockNetwork("127.0.0.1:"+port, registry)
	server, err := InitServer(node, network)
	assert.NoError(t, err)
	addr := "127.0.0.1:9997"
	registry.Register(addr)
	rpc := NewRPCMessage("FIND_VALUE", Payload{Key: "key", SourceContact: node.GetSelfContact()}, true)
	in := IncomingRPC{RPC: *rpc, Addr: addr}
	server.incoming <- in
	time.Sleep(500 * time.Millisecond)
	ch, ok := registry.Get(addr)
	assert.True(t, ok)
	select {
	case pkt := <-ch:
		var outRPC RPCMessage
		err := json.Unmarshal(pkt.data, &outRPC)
		assert.NoError(t, err)
		assert.Equal(t, "FIND_VALUE", outRPC.Type)
		assert.Nil(t, outRPC.Payload.Data)
	case <-time.After(1 * time.Second):
		t.Error("No FIND_VALUE response received")
	}
}

func Test_Server_ProcessRequest_Default_Error(t *testing.T) {
	port := "4324"
	node := &MockNodeAPI{Port: port}
	registry := NewMockRegistry()
	network := NewMockNetwork("127.0.0.1:"+port, registry)
	server, err := InitServer(node, network)
	assert.NoError(t, err)
	addr := "127.0.0.1:9996"
	registry.Register(addr)
	rpc := NewRPCMessage("UNKNOWN", Payload{SourceContact: node.GetSelfContact()}, true)
	in := IncomingRPC{RPC: *rpc, Addr: addr}
	server.incoming <- in
	time.Sleep(500 * time.Millisecond)
	ch, ok := registry.Get(addr)
	assert.True(t, ok)
	select {
	case pkt := <-ch:
		var outRPC RPCMessage
		err := json.Unmarshal(pkt.data, &outRPC)
		assert.NoError(t, err)
		assert.Equal(t, "ERROR", outRPC.Type)
		assert.Equal(t, node.GetSelfContact(), outRPC.Payload.TargetContact)
	case <-time.After(1 * time.Second):
		t.Error("No ERROR response received")
	}
}
