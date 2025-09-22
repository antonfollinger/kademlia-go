package kademlia

import (
	"log"
	"sort"
	"sync"
	"time"
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
	}

	node := &Node{
		Id:           kademliaID,
		RoutingTable: routingTable,
		Storage:      make(map[string][]byte),
	}

	return node, nil
}

// JoinNetwork performs an iterative lookup on the node's own ID to populate the routing table with nearby contacts
func (node *Node) JoinNetwork() error {

	selfID := node.GetSelfContact().ID
	bootstrapID := NewKademliaID("0000000000000000000000000000000000000000")

	if selfID.Equals(bootstrapID) {
		return nil
	} else {
		// If this is a peer (not bootstrap), ping the bootstrap node
		// Find bootstrap contact in routing table
		contacts := node.RoutingTable.FindClosestContacts(bootstrapID, 1)
		if len(contacts) > 0 {
			bootstrapContact := contacts[0]
			var err error
			for i := 0; i < 3; i++ { // Try up to 3 times
				_, err = node.Client.SendPingMessage(bootstrapContact)
				if err == nil {
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
			if err != nil {
				log.Printf("%s failed to ping bootstrap node on %s: %v\n", node.GetSelfContact().Address, bootstrapContact.Address, err)
			}
		} else {
			log.Printf("Bootstrap contact not found in routing table for %s\n", node.GetSelfContact().Address)
		}
		// Populate routing table with nearby contacts
		_, err := node.IterativeFindNode(node.Id)
		if err != nil {
			log.Printf("JoinNetwork: IterativeFindNode error: %v\n", err)
			return err
		}
		return nil
	}
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

	if contact.ID.Equals(node.GetSelfContact().ID) {
		return
	}

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
			//log.Printf("Contact %s responded to ping, not adding new contact %s\n", lastContact.String(), contact.String())
		} else {
			// Last contact did not respond, remove and add new contact
			bucket.list.Remove(lastContactElem)
			node.RoutingTable.AddContact(contact)
			//log.Printf("Contact %s did not respond, replaced with new contact %s\n", lastContact.String(), contact.String())
		}
	}
}

func (node *Node) LookupClosestContacts(target Contact) []Contact {
	return node.RoutingTable.FindClosestContacts(target.ID, alpha)
}

func (node *Node) IterativeFindNode(target *KademliaID) ([]Contact, error) {
	shortlist := node.LookupClosestContacts(NewContact(target, ""))
	if len(shortlist) == 0 {
		return nil, nil
	}
	queried := make(map[string]bool)
	inShortlist := make(map[string]bool)
	for _, c := range shortlist {
		if c.ID == nil {
			continue
		}
		inShortlist[c.ID.String()] = true
	}

	for {
		batch := []Contact{}
		for _, c := range shortlist {
			if c.ID == nil {
				continue
			}
			if !queried[c.ID.String()] && len(batch) < alpha {
				batch = append(batch, c)
			}
		}
		if len(batch) == 0 {
			break
		}
		results := make(chan []Contact, len(batch))
		for _, contact := range batch {
			if contact.ID == nil {
				results <- nil
				continue
			}
			queried[contact.ID.String()] = true
			go func(c Contact) {
				contacts, err := node.Client.SendFindNodeMessage(target, c)
				if err != nil || contacts == nil {
					results <- nil
					return
				}
				results <- contacts
			}(contact)
		}
		updated := false
		for i := 0; i < len(batch); i++ {
			var contacts []Contact
			select {
			case contacts = <-results:
				// got result
			case <-time.After(2 * time.Second):
				contacts = nil
			}
			if contacts == nil {
				continue
			}
			for _, c := range contacts {
				if c.ID == nil {
					continue
				}
				if !queried[c.ID.String()] && !inShortlist[c.ID.String()] {
					shortlist = append(shortlist, c)
					inShortlist[c.ID.String()] = true
					updated = true
				}
			}
		}
		if !updated {
			break
		}
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
	data, ok := node.Storage[hash]
	if !ok {
		return nil
	}
	return data
}

func (node *Node) Store(key string, data []byte) {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.Storage[key] = data
}

func (node *Node) PrintStore() {
	log.Printf("Store: %v\n", node.Storage)
}

func (node *Node) PrintRoutingTable() {
	printRoutingTable := func() {
		log.Printf("\n================= Routing Table %s =================\n", node.RoutingTable.me.ID.String())
		log.Printf("Self: %s (%s)\n", node.RoutingTable.me.Address, node.RoutingTable.me.ID.String())
		for i, bucket := range node.RoutingTable.buckets {
			if bucket.Len() == 0 {
				continue
			}
			log.Printf("Bucket %d:\n", i)
			for e := bucket.list.Front(); e != nil; e = e.Next() {
				contact := e.Value.(Contact)
				log.Printf("  - %s\t(%s)\t[%s]\n", contact.Address, contact.ID.String(), contact.distance.String())
			}
		}
		log.Println("==========================================================================================")
	}
	printRoutingTable()
}
