package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var myPort string

// MockNodeAPI for testing
type MockNodeAPI struct{}

func (m *MockNodeAPI) GetSelfContact() Contact {
	return Contact{ID: NewKademliaID("0000000000000000000000000000000000000001"), Address: "localhost:" + myPort}
}
func (m *MockNodeAPI) AddContact(contact Contact)                     {}
func (m *MockNodeAPI) LookupClosestContacts(target Contact) []Contact { return []Contact{} }
func (m *MockNodeAPI) IterativeFindNode(target *KademliaID) ([]Contact, error) {
	return []Contact{}, nil
}
func (m *MockNodeAPI) LookupData(hash string) []byte { return nil }
func (m *MockNodeAPI) Store(key string, data []byte) {}

func Test_Client_SendPingMessage_Timeout(t *testing.T) {
	myPort = "20001"
	// Use a real client but target an unreachable address
	//network, err := NewUDPNetwork(":" + myPort)
	//assert.NoError(t, err)
	client, err := InitClient(&MockNodeAPI{} /*network*/)
	assert.NoError(t, err)
	target := Contact{ID: NewKademliaID("0000000000000000000000000000000000000002"), Address: "127.0.0.1:65535"}
	resp, err := client.SendPingMessage(target)
	assert.Error(t, err)
	assert.Equal(t, RPCMessage{}, resp)
}

func Test_Client_SendMessage_Error(t *testing.T) {
	myPort = "20002"
	//network, err := NewUDPNetwork(":" + myPort)
	//assert.NoError(t, err)
	client, err := InitClient(&MockNodeAPI{} /*network*/)
	assert.NoError(t, err)
	// Invalid address
	target := Contact{ID: NewKademliaID("0000000000000000000000000000000000000003"), Address: "invalid:address"}
	msg := NewRPCMessage("PING", Payload{}, true)
	ch, err := client.SendMessage(target, msg)
	assert.Error(t, err)
	assert.Nil(t, ch)
}

func Test_Client_SendPingMessage_Timeout_Unreachable(t *testing.T) {
	myPort = "20003"
	//network, err := NewUDPNetwork(":" + myPort)
	//assert.NoError(t, err)
	client, err := InitClient(&MockNodeAPI{} /*network*/)
	assert.NoError(t, err)
	// Unreachable port
	target := Contact{ID: NewKademliaID("0000000000000000000000000000000000000004"), Address: "127.0.0.1:65534"}
	resp, err := client.SendPingMessage(target)
	assert.Error(t, err)
	assert.Equal(t, RPCMessage{}, resp)
}

func Test_Client_SendFindNodeMessage_Timeout_Unreachable(t *testing.T) {
	myPort = "20004"
	//network, err := NewUDPNetwork(":" + myPort)
	//assert.NoError(t, err)
	client, err := InitClient(&MockNodeAPI{} /*network*/)
	assert.NoError(t, err)
	// Unreachable port
	targetID := NewKademliaID("0000000000000000000000000000000000000005")
	contact := Contact{ID: NewKademliaID("0000000000000000000000000000000000000006"), Address: "127.0.0.1:65533"}
	contacts, err := client.SendFindNodeMessage(targetID, contact)
	assert.Error(t, err)
	assert.Nil(t, contacts)
}

func Test_Client_SendStoreMessage_Timeout_Unreachable(t *testing.T) {
	myPort = "20005"
	//network, err := NewUDPNetwork(":" + myPort)
	//assert.NoError(t, err)
	client, err := InitClient(&MockNodeAPI{} /*network*/)
	assert.NoError(t, err)
	// No reachable nodes, IterativeFindNode returns empty
	data := []byte("testdata")
	resp, err := client.SendStoreMessage(data)
	assert.Error(t, err)
	assert.Equal(t, RPCMessage{}, resp)
}

func Test_Client_SendFindValueMessage_Timeout_Unreachable(t *testing.T) {
	myPort = "20006"
	//network, err := NewUDPNetwork(":" + myPort)
	//assert.NoError(t, err)
	client, err := InitClient(&MockNodeAPI{} /*network*/)
	assert.NoError(t, err)
	// No reachable nodes, IterativeFindNode returns empty
	hash := "0000000000000000000000000000000000000007"
	resp, err := client.SendFindValueMessage(hash)
	assert.Error(t, err)
	assert.Equal(t, RPCMessage{}, resp)
}
