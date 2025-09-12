package kademlia

import (
	"fmt"
)

type Kademlia struct {
	RoutingTable *RoutingTable
	Network      *Network
	Storage      map[string][]byte
}

func InitKademlia(ip string, port int, bootstrap bool, bootstrapID string) *Kademlia {
	var kademliaID *KademliaID
	var contact Contact

	if bootstrap {
		kademliaID = NewKademliaID(bootstrapID)
		contact = NewContact(kademliaID, ip, port)
	} else {
		kademliaID = NewRandomKademliaID()
		contact = NewContact(kademliaID, ip, port)
	}
	routingTable := NewRoutingTable(contact)

	fmt.Printf("New node was created with: \n Address: %s\n Contact: %s\n ID: %s\n", contact.Address, contact.String(), contact.ID.String())

	return &Kademlia{
		RoutingTable: routingTable,
		Storage:      make(map[string][]byte),
	}
}

func (kademlia *Kademlia) SetNetworkInterface(network *Network) {
	kademlia.Network = network
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	return kademlia.RoutingTable.FindClosestContacts(target.ID, bucketSize)
}

func (kademlia *Kademlia) LookupData(hash string) []byte {
	return kademlia.Storage[hash]
}
func (kademlia *Kademlia) Store(key string, data []byte) {

	kademlia.Storage[key] = data
	fmt.Printf("Stored data with key %s\n", key)
}
