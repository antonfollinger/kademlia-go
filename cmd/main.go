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

	network, err := kademlia.Listen("0.0.0.0", port)
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}

	fmt.Printf("Node listening on UDP port %d\n", port)

	// Determine KademliaID
	if idStr := os.Getenv("KAD_ID"); idStr != "" {
		bootstrapnode := kademlia.InitKademlia(network.GetLocalIP(), port, true, os.Getenv("KAD_ID"))
		bootstrapnode.SetNetworkInterface(network)
		fmt.Println("Using fixed KademliaID:", idStr)

	}
	node := kademlia.InitKademlia(network.GetLocalIP(), port, false, "")
	node.SetNetworkInterface(network)

	peer := os.Getenv("PEER")
	if peer != "" {
		fmt.Println("Starting as peer, connecting to", peer)

		time.Sleep(2 * time.Second)

		parts := strings.Split(peer, ":")
		if len(parts) == 2 {
			bootstrapID := kademlia.NewKademliaID(os.Getenv("KAD_ID"))
			peerPort, _ := strconv.Atoi(parts[1])
			peerContact := kademlia.NewContact(bootstrapID, parts[0], peerPort) // << use same ID
			node.RoutingTable.AddContact(peerContact)

			fmt.Printf("Added peer contact: %+v, KademliaID: %s\n",
				peerContact, peerContact.ID.String())

			node.Network.SendPingMessage(&node.RoutingTable.FindClosestContacts(peerContact.ID, 1)[0])
			node.RoutingTable.Print()
		}
	} else {
		fmt.Println("Starting as bootstrap node")
	}

	select {}
}
