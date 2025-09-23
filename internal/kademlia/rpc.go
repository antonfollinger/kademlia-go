package kademlia

import (
	"github.com/google/uuid"
)

// Payload represents the data carried in an RPC message
type Payload struct {
	Contacts      []Contact `json:"contacts,omitempty"`
	SourceContact Contact   `json:"src_contact,omitempty"`
	TargetContact Contact   `json:"trgt_contact,omitempty"`
	Key           string    `json:"key,omitempty"`
	Data          []byte    `json:"data,omitempty"`
	Error         string    `json:"error,omitempty"`
}

// RPCMessage represents a message sent between nodes in the Kademlia network
type RPCMessage struct {
	Type     string  `json:"msg"`       // "PING", "STORE", "FIND_NODE", "FIND_VALUE"
	Payload  Payload `json:"payload"`   // The actual data being sent
	PacketID string  `json:"packet_id"` // Unique ID for the RPC call
	Query    bool    `json:"query"`     // Is this message a query (request) or a response
}

// NewRPCMessage creates a new RPCMessage with a unique PacketID
func NewRPCMessage(msgType string, payload Payload, query bool) *RPCMessage {
	newMessage := &RPCMessage{
		Type:     msgType,
		PacketID: uuid.New().String(),
		Payload:  payload,
		Query:    query,
	}
	return newMessage
}
