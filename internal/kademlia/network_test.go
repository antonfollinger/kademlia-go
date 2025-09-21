package kademlia

import (
	"fmt"
	"math/rand"
	"testing"
)

type MockNetwork struct {
	inbox     chan mockPacket
	dropRate  float64
	nodeAddrs map[string]chan mockPacket
}

func (m *MockNetwork) GetConn() string {
	return ""
}

type mockPacket struct {
	src  string
	data []byte
}

func (m *MockNetwork) SendMessage(addr string, data []byte) error {
	if rand.Float64() < m.dropRate {
		return fmt.Errorf("packet dropped")
	}
	if ch, ok := m.nodeAddrs[addr]; ok {
		ch <- mockPacket{src: addr, data: data}
		return nil
	}
	return fmt.Errorf("address not found")
}

func (m *MockNetwork) ReceiveMessage() (string, []byte, error) {
	pkt := <-m.inbox
	return pkt.src, pkt.data, nil
}

func (m *MockNetwork) Close() error {
	close(m.inbox)
	return nil
}

func Test_Kademlia_NetworkEmulation_WithPacketDrop(t *testing.T) {
	const nodeCount = 1000 // Change for testing
	const dropRate = 0.1   // Change for testing
	const messagesPerNode = 5

	addrMap := make(map[string]chan mockPacket, nodeCount)
	networks := make([]*MockNetwork, nodeCount)
	nodes := make([]*Node, nodeCount)
	clients := make([]*Client, nodeCount)
	servers := make([]*Server, nodeCount)

	// Setup address map and networks
	for i := 0; i < nodeCount; i++ {
		addr := fmt.Sprintf("node%d", i)
		inbox := make(chan mockPacket, 100)
		addrMap[addr] = inbox
	}

	// Create nodes, clients, and servers
	for i := 0; i < nodeCount; i++ {
		addr := fmt.Sprintf("node%d", i)
		networks[i] = &MockNetwork{
			inbox:     addrMap[addr],
			dropRate:  dropRate,
			nodeAddrs: addrMap,
		}
		node, err := InitNode(false, addr, "node0")
		if err != nil {
			t.Fatalf("Failed to init node %d: %v", i, err)
		}
		client, err := InitClient(node, networks[i])
		if err != nil {
			t.Fatalf("Failed to init client %d: %v", i, err)
		}
		server, err := InitServer(node, networks[i])
		if err != nil {
			t.Fatalf("Failed to init server %d: %v", i, err)
		}
		node.SetClient(client)
		nodes[i] = node
		clients[i] = client
		servers[i] = server
	}

	// Emulate sending PING messages between random nodes
	success := 0
	dropped := 0
	for i := 0; i < nodeCount; i++ {
		for j := 0; j < messagesPerNode; j++ {
			targetIdx := rand.Intn(nodeCount)
			if targetIdx == i {
				targetIdx = (targetIdx + 1) % nodeCount
			}
			targetContact := nodes[targetIdx].GetSelfContact()
			resp, err := clients[i].SendPingMessage(targetContact)
			if err != nil {
				dropped++
			} else if resp.Type == "PONG" {
				success++
			}
		}
	}

	t.Logf("\n\nTotal PINGs sent: %d \nSuccess: %d \nDropped: %d \nDropRate: %.2f\n\n", nodeCount*messagesPerNode, success, dropped, float64(dropped)/float64(nodeCount*messagesPerNode))
	if float64(dropped)/float64(nodeCount*messagesPerNode) < dropRate*0.8 || float64(dropped)/float64(nodeCount*messagesPerNode) > dropRate*1.2 {
		t.Errorf("Packet drop rate out of expected bounds")
	}
}
