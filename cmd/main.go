package main

import (
	"fmt"
	"os"
	"time"

	"github.com/antonfollinger/kademlia-go/tree/dev-RPC/internal/kademlia"
)

func main() {

	isBootstrap := os.Getenv("BOOTSTRAP")
	peer := os.Getenv("PEER")
	port := os.Getenv("PORT")

	var k *kademlia.Kademlia
	var err error

	if isBootstrap == "TRUE" {
		k, err = kademlia.InitKademlia(true, port)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize Kademlia: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Starting as bootstrap node")

		k.Server.RunServer()
	} else {
		k, err = kademlia.InitKademlia(false, port)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize Kademlia: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Starting as peer, connecting to", peer)

		time.Sleep(2 * time.Second)

		bootstrapID := kademlia.NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
		peerContact := kademlia.NewContact(bootstrapID, "0.0.0.0:1234")

		k.Node.RoutingTable.AddContact(peerContact)

		fmt.Printf("Added peer contact: %+v \n\n", peerContact)

		k.Server.RunServer()
		k.Client.SendPingMessage("0.0.0.0:1234")
	}

	select {}
}
