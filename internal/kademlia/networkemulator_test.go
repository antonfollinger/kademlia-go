package kademlia

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

type EmulatedNode struct {
	ID     string
	Node   *Node
	Client *Client
	Server *Server
}

type EmulatedNetwork struct {
	Nodes         map[string]*EmulatedNode
	DropRate      float64 // e.g. 0.05 for 5% drop
	mu            sync.Mutex
	DeliveredPkgs int
	DroppedPkgs   int
}

func NewEmulatedNetwork(numNodes int, dropRate float64) *EmulatedNetwork {
	net := &EmulatedNetwork{
		Nodes:    make(map[string]*EmulatedNode, numNodes+1),
		DropRate: dropRate,
	}

	// Create bootstrap node
	bootstrapID := "bootstrap"
	bootstrapAddr := "127.0.0.1:9000"
	bootstrapNode, _ := InitNode(true, bootstrapAddr, "")
	bootstrapClient, _ := InitClient(bootstrapNode)
	bootstrapServer, _ := InitServer(bootstrapNode)
	bootstrapNode.SetClient(bootstrapClient)
	net.Nodes[bootstrapID] = &EmulatedNode{ID: bootstrapID, Node: bootstrapNode, Client: bootstrapClient, Server: bootstrapServer}

	// Create peer nodes, passing bootstrap address
	for i := 0; i < numNodes; i++ {
		id := fmt.Sprintf("node-%d", i)
		node, _ := InitNode(false, id, bootstrapAddr)
		client, _ := InitClient(node)
		server, _ := InitServer(node)
		node.SetClient(client)
		net.Nodes[id] = &EmulatedNode{ID: id, Node: node, Client: client, Server: server}
	}

	// Ping bootstrap and join network for all peer nodes
	for i := 0; i < numNodes; i++ {
		id := fmt.Sprintf("node-%d", i)
		peer := net.Nodes[id]
		bootstrapContact := net.Nodes[bootstrapID].Node.GetSelfContact()
		_, _ = peer.Client.SendPingMessage(bootstrapContact)
		_ = peer.Node.JoinNetwork()
	}

	return net
}

// Emulate sending an RPC from src to dst, with packet drop
func (net *EmulatedNetwork) SendRPC(src, dst string, rpc *RPCMessage) error {
	net.mu.Lock()
	defer net.mu.Unlock()
	if rand.Float64() < net.DropRate {
		net.DroppedPkgs++
		return fmt.Errorf("packet dropped")
	}
	node, ok := net.Nodes[dst]
	if ok {
		// Directly process the RPC as if received by the server
		in := IncomingRPC{RPC: *rpc, Addr: nil}
		go node.Server.processRequest(in)
		net.DeliveredPkgs++
		return nil
	}
	return fmt.Errorf("destination not found")
}

func Test_EmulatedNetwork_KademliaStoreAndFind(t *testing.T) {
	numNodes := 100
	dropRate := 0.05 // 5% drop
	net := NewEmulatedNetwork(numNodes, dropRate)

	// Pick a source node and a key/data
	srcID := "node-0"
	key := "testkey"
	data := []byte("testdata")

	// Store data on k closest nodes using Kademlia logic
	srcNode := net.Nodes[srcID]
	resp, err := srcNode.Client.SendStoreMessage(data)
	if err != nil {
		t.Fatalf("StoreMessage failed: %v", err)
	}
	t.Logf("Store response: %+v", resp)

	// Try to find the value from another node
	dstID := "node-1"
	dstNode := net.Nodes[dstID]
	resp2, err := dstNode.Client.SendFindValueMessage(key)
	if err != nil {
		t.Fatalf("FindValueMessage failed: %v", err)
	}
	t.Logf("FindValue response: %+v", resp2)

	t.Logf("Delivered: %d, Dropped: %d, DropRate: %.2f", net.DeliveredPkgs, net.DroppedPkgs, dropRate)
}
