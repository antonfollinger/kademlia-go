package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/antonfollinger/kademlia_go/kademlia"
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

	peer := os.Getenv("PEER")
	if peer != "" {
		// peer format: nodeX:port
		time.Sleep(2 * time.Second)
		parts := strings.Split(peer, ":")
		if len(parts) == 2 {
			peerPort, _ := strconv.Atoi(parts[1])
			contact := &kademlia.Contact{Address: parts[0], Port: peerPort, ID: kademlia.NewRandomKademliaID()}
			fmt.Printf("Contact: %+v, KademliaID: %s\n", contact, contact.ID.String())
			network.SendPingMessage(contact)
		}
	}

	select {}
}
