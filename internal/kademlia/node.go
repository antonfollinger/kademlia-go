package kademlia

import (
	"fmt"
)

type Node struct {
	Id           *KademliaID
	RoutingTable *RoutingTable
	Storage      map[string][]byte
	Client       ClientAPI
}

type NodeAPI interface {
	GetSelfContact() Contact
	AddContact(contact Contact)
	LookupContact(target Contact) []Contact
	LookupData(hash string) []byte
	Store(key string, data []byte)
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

func (node *Node) SetClient(client ClientAPI) {
	node.Client = client
}

func (node *Node) GetSelfContact() (self Contact) {
	return node.RoutingTable.me
}

func (node *Node) AddContact(contact Contact) {

	bucket := node.RoutingTable.buckets[node.RoutingTable.getBucketIndex(contact.ID)]

	if bucket.Len() < bucketSize {
		node.RoutingTable.AddContact(contact)
	} else {
		// Ping last indexed contact
		lastContactElem := bucket.list.Back()
		lastContact := lastContactElem.Value.(Contact)
		resp, err := node.Client.SendPingMessage(lastContact)

		if err == nil && resp.Type == "PONG" {
			// Last contact responded, do not add new contact
			fmt.Printf("Contact %s responded to ping, not adding new contact %s\n", lastContact.String(), contact.String())
		} else {
			// Last contact did not respond, remove and add new contact
			bucket.list.Remove(lastContactElem)
			node.RoutingTable.AddContact(contact)
			fmt.Printf("Contact %s did not respond, replaced with new contact %s\n", lastContact.String(), contact.String())
		}
	}
}

func (node *Node) LookupContact(target Contact) []Contact {
	return node.RoutingTable.FindClosestContacts(target.ID, bucketSize)
}

func (node *Node) LookupData(hash string) []byte {
	return node.Storage[hash]
}

func (node *Node) Store(key string, data []byte) {
	node.Storage[key] = data
	fmt.Printf("Stored data with key %s\n", key)
}
