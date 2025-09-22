package kademlia

import (
	"fmt"
	"testing"
	"time"
)

const nodeCount = 50
const droprate = 0

func Test_Kademlia_NetworkEmulation_WithPacketDrop(t *testing.T) {

	// Create global MockRegistry
	registry := NewMockRegistry()

	totalNodes := nodeCount
	nodes := make([]*Kademlia, totalNodes)
	errs := make([]error, totalNodes)

	// Create bootstrap node first
	bootstrapAddr := "127.0.0.1:5000"
	bootstrapNode, bootstrapErr := InitKademlia(
		bootstrapAddr[len("127.0.0.1:"):],
		true,
		"",
		func(cfg *KademliaConfig) {
			cfg.isMockNetwork = true
			cfg.MockNetworkRegistry = registry
		},
	)
	nodes[0] = bootstrapNode
	errs[0] = bootstrapErr
	if bootstrapErr != nil {
		t.Fatalf("InitKademlia failed for bootstrap node: %v", bootstrapErr)
	}

	// Wait for bootstrap node to start listening
	time.Sleep(300 * time.Millisecond)

	// Create the rest of the nodes
	for i := 1; i < nodeCount; i++ {
		addr := fmt.Sprintf("127.0.0.1:%d", 5000+i)
		node, err := InitKademlia(
			addr[len("127.0.0.1:"):],
			false,
			bootstrapAddr,
			func(cfg *KademliaConfig) {
				cfg.isMockNetwork = true
				cfg.MockNetworkRegistry = registry
			},
		)
		nodes[i] = node
		errs[i] = err
		t.Logf("Created peer node %d as %v", i, node.Node.GetSelfContact())
	}
	for i, err := range errs {
		if err != nil {
			t.Fatalf("InitKademlia failed for node %d: %v", i, err)
		}
	}

	// Give goroutines time to start and stabilize
	time.Sleep(300 * time.Millisecond)

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

	// Randomly delete peer nodes from the network using droprate
	for i := 1; i < totalNodes; i++ {
		if randFloat := float32(i%10) / 10.0; randFloat < droprate {

			nodes[i].Client.Close()
			nodes[i].Server.Close()

			nodes[i] = nil
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
