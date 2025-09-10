package kademlia

import (
	"fmt"

	"github.com/google/uuid"
)

type Payload struct {
	Contacts      []Contact `json:"contacts,omitempty"`
	Contact       *Contact  `json:"contact,omitempty"`
	Key           string    `json:"key,omitempty"`
	Data          []byte    `json:"data,omitempty"`
	StringMessage string    `json:"string_message,omitempty"`
	Error         string    `json:"error,omitempty"`
}

type RPCMessage struct {
	Type     string  `json:"msg"`       // "PING", "STORE", "FIND_NODE", "FIND_VALUE"
	Payload  Payload `json:"payload"`   // The actual data being sent
	PacketID string  `json:"packet_id"` // Unique ID for the RPC call
	// SourceAddr *net.UDPAddr `json:"srcIP"`     // Source address of the message
	// DestAddr   *net.UDPAddr `json:"dstIP"`     // Destination address of the message
	// SourcePort int          `json:"srcPort"`   // Source port of the message
	// DestPort   int          `json:"dstPort"`   // Destination port of the message
	Query bool `json:"query"` // Is this message a query (request) or a response
}

func NewRPCMessage(msgType string, payload Payload /* srcAddr *net.UDPAddr, dstAddr *net.UDPAddr, */, query bool) *RPCMessage {
	newMessage := &RPCMessage{
		Type:     msgType,
		PacketID: uuid.New().String(),
		Payload:  payload,
		// SourceAddr: srcAddr,
		// DestAddr:   dstAddr,
		// SourcePort: srcAddr.Port,
		// DestPort:   dstAddr.Port,
		Query: query,
	}
	return newMessage
}

func (network *Network) SendPingMessage(contact *Contact) error {
	// Build payload with my own contact
	payload := Payload{
		Contact: &network.Kademlia.RoutingTable.me,
	}

	// Build RPCMessage with helper
	// srcAddr := &net.UDPAddr{IP: net.ParseIP(network.Kademlia.RoutingTable.me.Address), Port: network.Kademlia.RoutingTable.me.Port}
	// dstAddr := &net.UDPAddr{IP: net.ParseIP(contact.Address), Port: contact.Port}

	msg := NewRPCMessage("PING", payload /* srcAddr, dstAddr, */, true)

	network.SendMessage(contact, msg)
	return nil
}

func (network *Network) handleRPC(msg *RPCMessage) {
	switch msg.Type {
	case "PING":
		if msg.Query {

			network.Kademlia.RoutingTable.AddContact(*msg.Payload.Contact)
			replyPayload := Payload{Contact: &network.Kademlia.RoutingTable.me}
			reply := NewRPCMessage("PING", replyPayload, false)
			reply.PacketID = msg.PacketID // match request ID
			fmt.Printf("Got PING from %s:%d (ID=%s, PacketID=%s)\n",
				msg.Payload.Contact.Address,
				msg.Payload.Contact.Port,
				msg.Payload.Contact.ID.String(),
				msg.PacketID)
			_ = network.SendMessage(msg.Payload.Contact, reply)
		} else {
			if msg.Payload.Contact != nil {
				network.Kademlia.RoutingTable.AddContact(*msg.Payload.Contact)
				fmt.Printf("Got PONG from %s:%d (ID=%s, PacketID=%s)\n",
					msg.Payload.Contact.Address,
					msg.Payload.Contact.Port,
					msg.Payload.Contact.ID.String(),
					msg.PacketID)
			}
		}
	}
}
