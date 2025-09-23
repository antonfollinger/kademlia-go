package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockNodeAPI for testing
type MockNodeAPI struct {
	Port    string
	storage map[string][]byte
}

func (m *MockNodeAPI) GetSelfContact() Contact {
	return Contact{ID: NewKademliaID("0000000000000000000000000000000000000001"), Address: "localhost:" + m.Port}
}
func (m *MockNodeAPI) AddContact(contact Contact) {}
func (m *MockNodeAPI) LookupClosestContacts(target Contact) []Contact {
	return []Contact{m.GetSelfContact()}
}
func (m *MockNodeAPI) LookupData(hash string) []byte {
	if m.storage == nil {
		return nil
	}
	return m.storage[hash]
}
func (m *MockNodeAPI) Store(key string, data []byte) {
	if m.storage == nil {
		m.storage = make(map[string][]byte)
	}
	m.storage[key] = data
}
func (m *MockNodeAPI) IterativeFindNode(target *KademliaID) ([]Contact, error) {
	return []Contact{m.GetSelfContact()}, nil
}

func Test_Client_SendPingMessage_Timeout(t *testing.T) {
	port := "20001"
	registry := NewMockRegistry()
	network := NewMockNetwork("127.0.0.1:"+port, registry)
	client, err := InitClient(&MockNodeAPI{Port: port}, network)
	assert.NoError(t, err)
	target := Contact{ID: NewKademliaID("0000000000000000000000000000000000000002"), Address: "127.0.0.1:65535"}
	resp, err := client.SendPingMessage(target)
	assert.Error(t, err)
	assert.Equal(t, RPCMessage{}, resp)
}

func Test_Client_SendMessage_Error(t *testing.T) {
	port := "20002"
	registry := NewMockRegistry()
	network := NewMockNetwork("127.0.0.1:"+port, registry)
	client, err := InitClient(&MockNodeAPI{Port: port}, network)
	assert.NoError(t, err)
	// Invalid address
	target := Contact{ID: NewKademliaID("0000000000000000000000000000000000000003"), Address: "invalid:address"}
	msg := NewRPCMessage("PING", Payload{}, true)
	ch, err := client.SendMessage(target, msg)
	assert.Error(t, err)
	assert.Nil(t, ch)
}

func Test_Client_SendPingMessage_Timeout_Unreachable(t *testing.T) {
	port := "20003"
	registry := NewMockRegistry()
	network := NewMockNetwork("127.0.0.1:"+port, registry)
	client, err := InitClient(&MockNodeAPI{Port: port}, network)
	assert.NoError(t, err)
	// Unreachable port
	target := Contact{ID: NewKademliaID("0000000000000000000000000000000000000004"), Address: "127.0.0.1:65534"}
	resp, err := client.SendPingMessage(target)
	assert.Error(t, err)
	assert.Equal(t, RPCMessage{}, resp)
}

func Test_Client_SendFindNodeMessage_Timeout_Unreachable(t *testing.T) {
	port := "20004"
	registry := NewMockRegistry()
	network := NewMockNetwork("127.0.0.1:"+port, registry)
	client, err := InitClient(&MockNodeAPI{Port: port}, network)
	assert.NoError(t, err)
	// Unreachable port
	targetID := NewKademliaID("0000000000000000000000000000000000000005")
	contact := Contact{ID: NewKademliaID("0000000000000000000000000000000000000006"), Address: "127.0.0.1:65533"}
	contacts, err := client.SendFindNodeMessage(targetID, contact)
	assert.Error(t, err)
	assert.Nil(t, contacts)
}

func Test_Client_SendStoreMessage_Timeout_Unreachable(t *testing.T) {
	port := "20005"
	registry := NewMockRegistry()
	network := NewMockNetwork("127.0.0.1:"+port, registry)
	client, err := InitClient(&MockNodeAPI{Port: port}, network)
	assert.NoError(t, err)
	// No reachable nodes, IterativeFindNode returns empty
	data := []byte("testdata")
	resp, err := client.SendStoreMessage(data)
	assert.Error(t, err)
	assert.Equal(t, RPCMessage{}, resp)
}

func Test_Client_SendFindValueMessage_Timeout_Unreachable(t *testing.T) {
	port := "20006"
	registry := NewMockRegistry()
	network := NewMockNetwork("127.0.0.1:"+port, registry)
	client, err := InitClient(&MockNodeAPI{Port: port}, network)
	assert.NoError(t, err)
	// No reachable nodes, IterativeFindNode returns empty
	hash := "0000000000000000000000000000000000000007"
	resp, err := client.SendFindValueMessage(hash)
	assert.Error(t, err)
	assert.Equal(t, RPCMessage{}, resp)
}
