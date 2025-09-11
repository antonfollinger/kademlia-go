package kademlia

import (
	"github.com/google/uuid"
)

type Payload struct {
	Contacts      []Contact `json:"contacts,omitempty"`
	SourceContact *Contact  `json:"src_contact,omitempty"`
	Key           string    `json:"key,omitempty"`
	Data          []byte    `json:"data,omitempty"`
	Error         string    `json:"error,omitempty"`
}

type RPCMessage struct {
	Type     string  `json:"msg"`       // "PING", "STORE", "FIND_NODE", "FIND_VALUE", "OK"
	Payload  Payload `json:"payload"`   // The actual data being sent
	PacketID string  `json:"packet_id"` // Unique ID for the RPC call
}

func CreateRPCMessage(msgType string, payload Payload) *RPCMessage {
	newMessage := &RPCMessage{
		Type:     msgType,
		PacketID: uuid.New().String(),
		Payload:  payload,
	}
	return newMessage
}

/*
func ValidateRPCMessage(rpc *RPCMessage) bool {

}
*/
