package kademlia

import (
	"fmt"
	"sort"
	"sync"
)

const alpha = 3

type Node struct {
	Id           *KademliaID
	RoutingTable *RoutingTable
	Storage      map[string][]byte
	Client       ClientAPI
	mu           sync.Mutex
}

type NodeAPI interface {
	GetSelfContact() Contact
	AddContact(contact Contact)
	LookupClosestContacts(target Contact) []Contact
	IterativeFindNode(target *KademliaID) ([]Contact, error)
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
	node.mu.Lock()
	defer node.mu.Unlock()
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

func (node *Node) LookupClosestContacts(target Contact) []Contact {
	return node.RoutingTable.FindClosestContacts(target.ID, alpha)
}

func (node *Node) IterativeFindNode(target *KademliaID) ([]Contact, error) {
	shortlist := node.LookupClosestContacts(NewContact(target, ""))
	queried := make(map[string]bool)

	for {
		batch := []Contact{}

		// take alpha closest unqueried nodes
		for _, c := range shortlist {
			if !queried[c.ID.String()] && len(batch) < alpha {
				batch = append(batch, c)
			}
		}

		// if we cannot find any closer nodes then we are done
		if len(batch) == 0 {
			break
		}

		results := make(chan []Contact, len(batch))
		for _, contact := range batch {
			queried[contact.ID.String()] = true
			go func(c Contact) {
				contacts, err := node.Client.SendFindNodeMessage(target, c)
				if err != nil {
					results <- nil
					return
				}
				results <- contacts
			}(contact)
		}

		updated := false
		for i := 0; i < len(batch); i++ {
			contacts := <-results
			for _, c := range contacts {
				if !queried[c.ID.String()] {
					shortlist = append(shortlist, c)
					updated = true
				}
			}
		}

		if !updated {
			break
		}

		// Sort the shortlist slice
		sort.Slice(shortlist, func(i, j int) bool {
			di := shortlist[i].ID.CalcDistance(target)
			dj := shortlist[j].ID.CalcDistance(target)
			return di.Less(dj)
		})
	}

	k := alpha
	if len(shortlist) < k {
		k = len(shortlist)
	}
	return shortlist[:k], nil
}

func (node *Node) LookupData(hash string) []byte {
	return node.Storage[hash]
}

func (node *Node) Store(key string, data []byte) {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.Storage[key] = data
	fmt.Printf("Stored data with key %s\n", key)
}
