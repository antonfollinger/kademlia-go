package kademlia

import (
	"fmt"
	"sync"
	"testing"
)

func Test_Kademlia_NetworkEmulation_WithPacketDrop(t *testing.T) {
	const nodeCount = 100
	const droprate = 0.5

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

	// Bootstrap node puts up a single value
	bootstrap := nodes[0]
	value := "test123"
	// Use client to store value in the network
	req, err := bootstrap.Client.SendStoreMessage([]byte(value))
	// Hashed key for value
	key := req.Payload.Key
	if err != nil {
		t.Fatalf("Bootstrap node failed to store value: %v", err)
	} else {
		t.Logf("Bootstrap node stored value '%s' under key '%s'", value, key)
	}

	// Randomly delete nodes from the network using droprate
	deleted := make(map[int]bool)
	for i := 1; i < totalNodes; i++ {
		// Don't delete bootstrap node
		if randFloat := float32(i%10) / 10.0; randFloat < droprate {
			nodes[i] = nil
			deleted[i] = true
			t.Logf("Node %d deleted from network", i)
		}
	}

	// Bootstrap tries to fetch the value again
	res, err := bootstrap.Client.SendFindValueMessage(key)
	val := string(res.Payload.Data)
	if err != nil {
		t.Errorf("Bootstrap node failed to fetch value after deletions: %v", err)
	} else if val != value {
		t.Errorf("Bootstrap node fetched wrong value: got '%v', want '%v'", val, value)
	} else {
		t.Errorf("Bootstrap node successfully fetched value '%s' after deletions", val)
	}
}
