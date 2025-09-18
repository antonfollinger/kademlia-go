package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var myPort = "20001"

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
	// Use a real client but target an unreachable address
	client, err := InitClient(&MockNodeAPI{})
	assert.NoError(t, err)
	target := Contact{ID: NewKademliaID("0000000000000000000000000000000000000002"), Address: "127.0.0.1:65535"}
	resp, err := client.SendPingMessage(target)
	assert.Error(t, err)
	assert.Equal(t, RPCMessage{}, resp)
}

func Test_Client_SendMessage_Error(t *testing.T) {
	client, err := InitClient(&MockNodeAPI{})
	assert.NoError(t, err)
	// Invalid address
	target := Contact{ID: NewKademliaID("0000000000000000000000000000000000000003"), Address: "invalid:address"}
	msg := NewRPCMessage("PING", Payload{}, true)
	ch, err := client.SendMessage(target, msg)
	assert.Error(t, err)
	assert.Nil(t, ch)
}
