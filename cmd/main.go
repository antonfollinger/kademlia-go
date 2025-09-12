package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/antonfollinger/kademlia_go/internal/kademlia"
)

func parsePeersEnv() []kademlia.Contact {
	var out []kademlia.Contact
	val := os.Getenv("PEERS")
	if val == "" {
		val = os.Getenv("PEER")
	}
	if val == "" {
		return out
	}
	items := strings.Split(val, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		hostPort, idStr, found := strings.Cut(item, "@")
		if !found {
			idStr = "0000000000000000000000000000000000000000"
		}
		host, portStr, ok := strings.Cut(hostPort, ":")
		if !ok {
			continue
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}
		id := kademlia.NewKademliaID(idStr)
		c := kademlia.NewContact(id, host, port)
		out = append(out, c)
	}
	return out
}

func main() {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "9001"
	}
	port, _ := strconv.Atoi(portStr)

	// Start UDP network
	network, err := kademlia.Listen("0.0.0.0", port)
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}
	fmt.Printf("Node listening on UDP port %d\n", port)

	// Create the Kademlia node
	var node *kademlia.Kademlia
	idStr := os.Getenv("KAD_ID")
	if idStr != "" {
		node = kademlia.InitKademlia(network.GetLocalIP(), port, true, idStr)
		fmt.Println("Using fixed KademliaID:", idStr)
	} else {
		node = kademlia.InitKademlia(network.GetLocalIP(), port, false, "")
	}
	node.SetNetworkInterface(network)
	network.Kademlia = node

	initialPeers := parsePeersEnv()

	if len(initialPeers) > 0 {
		fmt.Printf("Initial peers: %d\n", len(initialPeers))
	}

	if len(initialPeers) > 0 {
		// === PEER NODES ===
		go func() {
			me := kademlia.GetMe(node.RoutingTable)

			// Step 1: PING loop for a while
			for i := 0; i < 3; i++ { // do 3 rounds of PINGs
				for _, c := range initialPeers {
					if c.ID.Equals(me.ID) {
						continue
					}
					ping := kademlia.NewRPCMessage("PING", kademlia.Payload{
						TargetContact: &c,
						SourceContact: &me,
					}, true)
					_ = node.Network.SendMessage(&c, ping)
					fmt.Printf("Sent PING to %s:%d (ID=%s)\n", c.Address, c.Port, c.ID)
				}
				time.Sleep(10 * time.Second)
			}

			// Step 2: After delay, send a single FIND_NODE(self) to bootstrap
			time.Sleep(5 * time.Second)
			for _, c := range initialPeers {
				findNode := kademlia.NewRPCMessage("FIND_NODE", kademlia.Payload{
					TargetContact: &me,
					SourceContact: &me,
				}, true)
				_ = node.Network.SendMessage(&c, findNode)
				fmt.Printf("Sent FIND_NODE(self) to %s:%d (ID=%s)\n", c.Address, c.Port, c.ID)
			}

			// Step 3: Wait a bit for response, then dump routing table once
			time.Sleep(5 * time.Second)
			fmt.Println("=== Final Routing Table ===")
			node.RoutingTable.Print()
		}()
	} else {
		// === BOOTSTRAP NODE ===
		fmt.Println("Running as bootstrap (will just collect PINGs)")
		// Bootstrap can print its routing table once after some time
		go func() {
			time.Sleep(40 * time.Second)
			fmt.Println("=== Bootstrap Routing Table ===")
			node.RoutingTable.Print()
		}()
	}

	select {} // keep running
}
