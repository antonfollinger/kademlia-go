package kademlia

import (
	"fmt"
)

type Kademlia struct {
	RoutingTable *RoutingTable
	Network      *Network
	Storage      map[string][]byte
}

func InitKademlia(ip string, port int) *Kademlia {
	kademliaID := NewRandomKademliaID()
	contact := NewContact(kademliaID, ip, port)
	routingTable := NewRoutingTable(contact)

	fmt.Printf("New node was created with: \n Address: %s\n Contact: %s\n ID: %s\n", contact.Address, contact.String(), contact.ID.String())

	return &Kademlia{
		RoutingTable: routingTable,
		Storage:      make(map[string][]byte),
	}
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	return kademlia.RoutingTable.FindClosestContacts(target.ID, bucketSize)
}

func (kademlia *Kademlia) LookupData(hash string) {
	if value, exists := kademlia.Storage[hash]; exists {
		fmt.Printf("Data found locally: %s\n", string(value))
		kademlia.Network.SendFindDataMessage(string(value))
		return
	}
}
func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
