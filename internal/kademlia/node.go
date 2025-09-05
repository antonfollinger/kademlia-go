package kademlia

type Node struct {
	ID           *NodeID
	Address      string
	RoutingTable RoutingTable
	DataStore    *DataStore
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
