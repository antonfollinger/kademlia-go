package kademlia

type Node struct {
	id           *KademliaID
	routingTable *RoutingTable
	storage      map[string][]byte
}

func InitNode(isBootstrap bool, ip string) (*Node, error) {

	// NodeID
	var nodeID *KademliaID
	if isBootstrap {
		nodeID = NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	} else {
		nodeID = NewRandomKademliaID()
	}

	c := NewContact(nodeID, ip)
	rt := NewRoutingTable(c)

	n := &Node{
		id:           nodeID,
		routingTable: rt,
		storage:      make(map[string][]byte),
	}

	return n, nil

}

func (node *Node) LookupContact(target *Contact) {
	// TODO
}

func (node *Node) LookupData(hash string) {
	// TODO
}

func (node *Node) Store(data []byte) {
	// TODO
}
