package kademlia

import (
	"fmt"
)

type Node struct {
	Id           *KademliaID
	RoutingTable *RoutingTable
	Storage      map[string][]byte
}

type NodeAPI interface {
	LookupContact(target *Contact)
	LookupData()
	Store()
}

func InitNode(isBootstrap bool, ip string, bootstrapIP string) (*Node, error) {

	var kademliaID *KademliaID
	var me Contact

	// Create routing table
	if isBootstrap {
		kademliaID = NewKademliaID("0000000000000000000000000000000000000000")
	} else {
		kademliaID = NewRandomKademliaID()
	}
	me = NewContact(kademliaID, ip)
	routingTable := NewRoutingTable(me)

	// Create and add bootstrap contact if node is a peer
	if !isBootstrap {
		bootstrap := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), bootstrapIP)
		routingTable.AddContact(bootstrap)
		fmt.Printf("\nBootstrap added with: \n Address: %s\n Contact: %s\n ID: %s\n", bootstrap.Address, bootstrap.String(), bootstrap.ID.String())
	}

	fmt.Printf("\nNew node was created with: \n Address: %s\n Contact: %s\n ID: %s\n\n", me.Address, me.String(), me.ID.String())

	node := &Node{
		Id:           kademliaID,
		RoutingTable: routingTable,
		Storage:      make(map[string][]byte),
	}

	return node, nil
}

func (node *Node) LookupContact(target *Contact) []Contact {
	return node.RoutingTable.FindClosestContacts(target.ID, bucketSize)
}

func (node *Node) LookupData(hash string) []byte {
	return node.Storage[hash]
}

func (node *Node) Store(key string, data []byte) {

	node.Storage[key] = data
	fmt.Printf("Stored data with key %s\n", key)
}
