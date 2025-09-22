package kademlia

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func Test_Kademlia_NetworkEmulation_WithPacketDrop(t *testing.T) {
	const nodeCount = 100
	const messagesPerNode = 1

	// Create global MockRegistry
	registry := NewMockRegistry()

	// Create all nodes (bootstrap first, then peers) in parallel
	totalNodes := nodeCount
	nodes := make([]*Kademlia, totalNodes)
	errs := make([]error, totalNodes)

	var wg sync.WaitGroup
	wg.Add(totalNodes)

	for i := 0; i < totalNodes; i++ {
		go func(i int) {
			defer wg.Done()
			var addr string
			var isBootstrap bool
			var bootstrapIP string
			if i == 0 {
				// Bootstrap node
				addr = "127.0.0.1:1234"
				isBootstrap = true
				bootstrapIP = ""
			} else {
				// Peer node
				addr = fmt.Sprintf("127.0.0.1:%d", 5000+i)
				isBootstrap = false
				bootstrapIP = "127.0.0.1:1234"
			}
			node, err := InitKademlia(addr[len("127.0.0.1:"):], isBootstrap, bootstrapIP, func(cfg *KademliaConfig) {
				cfg.isMockNetwork = true
				cfg.MockNetworkRegistry = registry
				cfg.DropRate = 0.1
				cfg.SkipBootstrapPing = true
			})
			if err != nil {
				t.Logf("Node %d creation failed: %v", i, err)
			} else {
				t.Logf("Node %d created successfully", i)
			}
			nodes[i] = node
			errs[i] = err
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Fatalf("InitKademlia failed for node %d: %v", i, err)
		}
	}

	var success, dropped int64
	var sendWg sync.WaitGroup
	sendWg.Add(totalNodes - 1)
	for i := 1; i < totalNodes; i++ {
		go func(i int) {
			defer sendWg.Done()
			sender := nodes[i]
			receiver := nodes[i-1]
			t.Logf("Node %d sending %d PINGs to node %d", i, messagesPerNode, i-1)
			for m := 0; m < messagesPerNode; m++ {
				_, err := sender.Client.SendPingMessage(receiver.Node.GetSelfContact())
				if err != nil {
					t.Logf("PING from node %d to node %d failed: %v", i, i-1, err)
					atomic.AddInt64(&dropped, 1)
				} else {
					t.Logf("PING from node %d to node %d succeeded", i, i-1)
					atomic.AddInt64(&success, 1)
				}
			}
		}(i)
	}
	sendWg.Wait()

	t.Logf("Total PINGs sent: %d, Success: %d, Dropped: %d", nodeCount*messagesPerNode, success, dropped)
	if dropped == 0 {
		t.Errorf("No packets were dropped, expected some drops in mock network")
	}
}
