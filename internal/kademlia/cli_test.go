package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockClient for CLI tests
type MockClientCLI struct{}

func (mc *MockClientCLI) SendPingMessage(target Contact) (RPCMessage, error) {
	return RPCMessage{Type: "PONG"}, nil
}
func (mc *MockClientCLI) SendFindNodeMessage(target *KademliaID, contact Contact) ([]Contact, error) {
	return []Contact{}, nil
}
func (mc *MockClientCLI) SendStoreMessage(data []byte) (RPCMessage, error) {
	return RPCMessage{Payload: Payload{Key: "testhash"}, PacketID: "packet123"}, nil
}
func (mc *MockClientCLI) SendFindValueMessage(hash string) (RPCMessage, error) {
	return RPCMessage{
		Payload: Payload{
			Key:           hash,
			Data:          []byte("testdata"),
			SourceContact: Contact{ID: NewKademliaID("1234567891234567891234567891234567891234"), Address: "addr"},
		},
	}, nil
}

func Test_Node_Put_Success(t *testing.T) {
	node, _ := InitNode(true, "localhost:9000", "")
	node.SetClient(&MockClientCLI{})
	result, err := node.Put("somedata")
	assert.NoError(t, err)
	assert.Contains(t, result, "✅ Content stored!")
	assert.Contains(t, result, "Hash: testhash")
	assert.Contains(t, result, "Packet ID: packet123")
}

func Test_Node_Put_Error(t *testing.T) {
	node, _ := InitNode(true, "localhost:9001", "")
	node.SetClient(&MockClientError{})
	result, err := node.Put("faildata")
	assert.Error(t, err)
	assert.Empty(t, result)
}

type MockClientError struct{}

func (mc *MockClientError) SendPingMessage(target Contact) (RPCMessage, error) {
	return RPCMessage{}, nil
}
func (mc *MockClientError) SendFindNodeMessage(target *KademliaID, contact Contact) ([]Contact, error) {
	return nil, nil
}
func (mc *MockClientError) SendStoreMessage(data []byte) (RPCMessage, error) {
	return RPCMessage{}, assert.AnError
}
func (mc *MockClientError) SendFindValueMessage(hash string) (RPCMessage, error) {
	return RPCMessage{}, nil
}

func Test_Node_Get_Success(t *testing.T) {
	node, _ := InitNode(true, "localhost:9002", "")
	node.SetClient(&MockClientCLI{})
	result, err := node.Get("testhash")
	assert.NoError(t, err)
	assert.Contains(t, result, "✅ Content retrieved!")
	assert.Contains(t, result, "Hash: testhash")
	assert.Contains(t, result, "Content: testdata")
	assert.Contains(t, result, "Source: 1234567891234567891234567891234567891234")
}

func Test_Node_Get_Error(t *testing.T) {
	node, _ := InitNode(true, "localhost:9003", "")
	node.SetClient(&MockClientError{})
	result, err := node.Get("failhash")
	assert.Error(t, err)
	assert.Empty(t, result)
}
