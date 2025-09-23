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
			for range 3 { // Try up to 3 times
				rpc, err := node.Client.SendPingMessage(bootstrapContact)
				if err == nil {
					node.AddContact(rpc.Payload.SourceContact)
					break
				}
				time.Sleep(300 * time.Millisecond)
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

func (n *Node) AddContact(c Contact) {

	// Should not add self
	if c == n.GetSelfContact() {
		return
	}

	n.mu.Lock()
	bucketIndex := n.RoutingTable.getBucketIndex(c.ID)
	bucket := n.RoutingTable.buckets[bucketIndex]
	// Calculate and set the contact's distance field before adding
	c.distance = n.Id.CalcDistance(c.ID)
	if bucket.Len() < bucketSize {
		bucket.AddContact(c)
		n.mu.Unlock()
		return
	}
	// If full, copy LRU contact to ping after releasing lock
	lru := bucket.list.Back().Value.(Contact)
	n.mu.Unlock()

	// Ping LRU outside lock
	alive := false
	if resp, err := n.Client.SendPingMessage(lru); err == nil && resp.Type == "PONG" {
		alive = true
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	bucketIndex = n.RoutingTable.getBucketIndex(c.ID)
	bucket = n.RoutingTable.buckets[bucketIndex] // re-fetch in case table changed
	if bucket.Len() < bucketSize {
		bucket.AddContact(c)
	} else if !alive {
		// evict old LRU and add new contact
		bucket.list.Remove(bucket.list.Back())
		bucket.AddContact(c)
	}
	// else: keep existing LRU, drop new
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
