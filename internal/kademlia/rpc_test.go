package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewRPCMessage_Creation(t *testing.T) {
	payload := Payload{
		Key:   "testkey",
		Data:  []byte("testdata"),
		Error: "",
	}
	msg := NewRPCMessage("STORE", payload, true)
	assert.Equal(t, "STORE", msg.Type)
	assert.Equal(t, payload, msg.Payload)
	assert.True(t, msg.Query)
	assert.NotEmpty(t, msg.PacketID)
}

func Test_RPCMessage_Contacts(t *testing.T) {
	contact1 := Contact{ID: NewKademliaID("0000000000000000000000000000000000000001"), Address: "localhost:8001"}
	contact2 := Contact{ID: NewKademliaID("0000000000000000000000000000000000000002"), Address: "localhost:8002"}
	payload := Payload{
		Contacts:      []Contact{contact1, contact2},
		SourceContact: contact1,
		TargetContact: contact2,
	}
	msg := NewRPCMessage("FIND_NODE", payload, true)
	assert.Equal(t, 2, len(msg.Payload.Contacts))
	assert.Equal(t, contact1, msg.Payload.SourceContact)
	assert.Equal(t, contact2, msg.Payload.TargetContact)
}

func Test_RPCMessage_ErrorField(t *testing.T) {
	payload := Payload{
		Error: "something went wrong",
	}
	msg := NewRPCMessage("PING", payload, false)
	assert.Equal(t, "something went wrong", msg.Payload.Error)
	assert.False(t, msg.Query)
}
