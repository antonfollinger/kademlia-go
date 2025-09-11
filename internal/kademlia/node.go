package kademlia

const (
	alpha = 3
)

type Node struct {
	Id           *KademliaID
	RoutingTable *RoutingTable
	Storage      map[string][]byte
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
		Id:           nodeID,
		RoutingTable: rt,
		Storage:      make(map[string][]byte),
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
