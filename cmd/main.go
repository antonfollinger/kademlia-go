package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/antonfollinger/kademlia_go/internal/kademlia"
)

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
	peer := os.Getenv("PEER")
	if idStr != "" && peer == "" {
		node = kademlia.InitKademlia(network.GetLocalIP(), port, true, idStr)
		fmt.Println("Using fixed KademliaID:", idStr)
	} else {
		node = kademlia.InitKademlia(network.GetLocalIP(), port, false, "")
	}
	node.SetNetworkInterface(network)
	network.Kademlia = node // link back so RPC handler can see it

	if peer != "" {
		// --- run as peer ---
		fmt.Println("Starting as peer, connecting to", peer)
		time.Sleep(2 * time.Second) // give bootstrap time to start

		parts := strings.Split(peer, ":")
		if len(parts) == 2 {
			peerID := kademlia.NewKademliaID(os.Getenv("KAD_ID"))
			peerPort, _ := strconv.Atoi(parts[1])
			peerContact := kademlia.NewContact(peerID, parts[0], peerPort)

			// Add bootstrap to routing table
			node.RoutingTable.AddContact(peerContact)
			fmt.Printf("Added peer contact: %s (ID=%s)\n",
				peerContact.Address, peerContact.ID.String())

			// Send PING to bootstrap
			if err := node.Network.SendPingMessage(&peerContact); err != nil {
				fmt.Println("Error sending PING:", err)
			}
		}
	} else {
		// --- run as bootstrap ---
		fmt.Println("Starting as bootstrap node")

		// Periodically print routing table
		go func() {
			for {
				time.Sleep(10 * time.Second)
				fmt.Println("=== Routing Table Dump ===")
				node.RoutingTable.Print()
			}
		}()
	}

	select {} // keep running
}
