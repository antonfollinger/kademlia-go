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

func Test_MockNetwork_PacketDrop_Emulation(t *testing.T) {
	const nodeCount = 1000
	const dropRate = 0.1
	const messagesPerNode = 10

	// Seed random for reproducibility
	rand.Seed(1)

	// Create address map and nodes
	addrMap := make(map[string]chan mockPacket, nodeCount)
	nodes := make([]*MockNetwork, nodeCount)
	for i := 0; i < nodeCount; i++ {
		addr := fmt.Sprintf("node%d", i)
		inbox := make(chan mockPacket, 100)
		addrMap[addr] = inbox
		nodes[i] = &MockNetwork{
			inbox:     inbox,
			dropRate:  dropRate,
			nodeAddrs: addrMap,
		}
	}

	// Send messages between random nodes
	success := 0
	dropped := 0
	for i := 0; i < nodeCount; i++ {
		for j := 0; j < messagesPerNode; j++ {
			target := fmt.Sprintf("node%d", rand.Intn(nodeCount))
			data := []byte(fmt.Sprintf("msg from node%d", i))
			err := nodes[i].SendMessage(target, data)
			if err != nil {
				dropped++
			} else {
				success++
			}
		}
	}

	t.Logf("Total sent: %d, Success: %d, Dropped: %d, DropRate: %.2f", nodeCount*messagesPerNode, success, dropped, float64(dropped)/float64(nodeCount*messagesPerNode))
	if float64(dropped)/float64(nodeCount*messagesPerNode) < dropRate*0.8 || float64(dropped)/float64(nodeCount*messagesPerNode) > dropRate*1.2 {
		t.Errorf("Packet drop rate out of expected bounds")
	}
}
