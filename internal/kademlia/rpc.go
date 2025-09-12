package kademlia

import (
	"fmt"

	"github.com/google/uuid"
)

type Payload struct {
	Contacts      []Contact `json:"contacts,omitempty"`
	SourceContact *Contact  `json:"src_contact,omitempty"`
	TargetContact *Contact  `json:"dst_contact,omitempty"`
	Key           string    `json:"key,omitempty"`
	Data          []byte    `json:"data,omitempty"`
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

func (network *Network) SendPingMessage(msg *RPCMessage) error {
	// Build payload with my own contact
	payload := Payload{
		SourceContact: &network.Kademlia.RoutingTable.me,
	}

	if msg.Query {
		if msg.Payload.SourceContact != nil {
			network.Kademlia.RoutingTable.AddContact(*msg.Payload.SourceContact)
			reply := NewRPCMessage("PING", payload, false)
			reply.PacketID = msg.PacketID // match request ID
			fmt.Printf("Got PING from %s:%d (ID=%s, PacketID=%s)\n",
				msg.Payload.SourceContact.Address,
				msg.Payload.SourceContact.Port,
				msg.Payload.SourceContact.ID.String(),
				msg.PacketID)
			_ = network.SendMessage(msg.Payload.SourceContact, reply)
		}
	} else {
		if msg.Payload.SourceContact != nil {
			network.Kademlia.RoutingTable.AddContact(*msg.Payload.SourceContact)
			fmt.Printf("Got PONG from %s:%d (ID=%s, PacketID=%s)\n",
				msg.Payload.SourceContact.Address,
				msg.Payload.SourceContact.Port,
				msg.Payload.SourceContact.ID.String(),
				msg.PacketID)
		}
	}

	return nil
}

func (network *Network) SendFindContactMessage(msg *RPCMessage) {
	if msg.Query {
		contacts := network.Kademlia.RoutingTable.FindClosestContacts(msg.Payload.TargetContact.ID, 8)

		payload := Payload{
			Contacts:      contacts,
			SourceContact: &network.Kademlia.RoutingTable.me,
		}
		reply := NewRPCMessage("FIND_NODE", payload, false)
		reply.PacketID = msg.PacketID // preserve request ID

		if msg.Payload.SourceContact != nil {
			_ = network.SendMessage(msg.Payload.SourceContact, reply)
		}
	} else {
		if msg.Payload.SourceContact != nil {
			network.Kademlia.RoutingTable.AddContact(*msg.Payload.SourceContact)
		}

		if msg.Payload.Contacts != nil {
			fmt.Printf("Got FIND_NODE response with %d contacts (PacketID=%s)\n",
				len(msg.Payload.Contacts), msg.PacketID)
			for _, contact := range msg.Payload.Contacts {
				network.Kademlia.RoutingTable.AddContact(contact)
			}
		}
	}
}

func (network *Network) SendStoreMessage(msg *RPCMessage) {
	if msg.Query {
		network.Kademlia.Store(msg.Payload.Key, msg.Payload.Data)

		payload := Payload{
			SourceContact: &network.Kademlia.RoutingTable.me,
		}
		msgout := NewRPCMessage("STORE", payload, false)
		network.SendMessage(msg.Payload.SourceContact, msgout)

	} else {
		if msg.Payload.SourceContact != nil {
			network.Kademlia.RoutingTable.AddContact(*msg.Payload.SourceContact)

			network.Kademlia.RoutingTable.AddContact(*msg.Payload.SourceContact)
			fmt.Printf("Got STORE ACK from %s:%d (ID=%s, PacketID=%s)\n",
				msg.Payload.SourceContact.Address,
				msg.Payload.SourceContact.Port,
				msg.Payload.SourceContact.ID.String(),
				msg.PacketID)
			network.Kademlia.RoutingTable.AddContact(*msg.Payload.SourceContact)
		}
	}

}

func (network *Network) SendFindValueMessage(msg *RPCMessage) {
	if msg.Query {
		payload := Payload{
			Data:          network.Kademlia.LookupData(msg.Payload.Key),
			SourceContact: &network.Kademlia.RoutingTable.me,
		}
		msgout := NewRPCMessage("FIND_VALUE", payload, false)
		network.SendMessage(msg.Payload.SourceContact, msgout)
	} else {
		network.Kademlia.RoutingTable.AddContact(*msg.Payload.SourceContact)

		if msg.Payload.Data != nil {
			fmt.Printf("Got FIND_VALUE response with data: %s (PacketID=%s)\n",
				string(msg.Payload.Data),
				msg.PacketID)
		}
	}
}

func (network *Network) handleRPC(msg *RPCMessage) {
	switch msg.Type {
	case "PING":
		network.SendPingMessage(msg)
	case "FIND_NODE":
		network.SendFindContactMessage(msg)
	case "STORE":
		network.SendStoreMessage(msg)
	case "FIND_VALUE":
		network.SendFindValueMessage(msg)
	default:
		fmt.Println("Unknown RPC message type:", msg.Type)
	}

}
