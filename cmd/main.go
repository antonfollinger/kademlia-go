package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/antonfollinger/kademlia_go/internal/kademlia"
)

func getContainerIP() string {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return ""
}

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

	// Use InitKademlia to initialize the node
	node := kademlia.InitKademlia(getContainerIP(), port)
	node.Network = network // set the network field if needed

	peer := os.Getenv("PEER")
	if peer != "" {
		// peer format: ip:port
		time.Sleep(2 * time.Second)
		parts := strings.Split(peer, ":")
		if len(parts) == 2 {
			peerPort, _ := strconv.Atoi(parts[1])
			peerContact := kademlia.NewContact(kademlia.NewRandomKademliaID(), parts[0], peerPort)
			node.RoutingTable.AddContact(peerContact)
			fmt.Printf("Peer Contact: %+v, KademliaID: %s\n", peerContact, peerContact.ID.String())
			network.SendPingMessage(&node.RoutingTable.FindClosestContacts(peerContact.ID, 1)[0])
			node.RoutingTable.Print()
		}
	}

	select {}
}
