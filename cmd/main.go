package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/antonfollinger/kademlia_go/internal/kademlia"
)

func main() {
	isBootstrap := os.Getenv("ISBOOTSTRAP")
	port := os.Getenv("PORT")

	var k *kademlia.Kademlia
	var kadErr error
	var bootstrapIP string

	if isBootstrap == "TRUE" {
		k, kadErr = kademlia.InitKademlia(port, true, "")
		if kadErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize Kademlia: %v\n", kadErr)
			os.Exit(1)
		}

	} else {
		// Bootstrap Node info with retry
		bootstrapNode := os.Getenv("BOOTSTRAPNODE")
		var bootStrapAddr []net.IP
		var IPerr error
		for i := 0; i < 5; i++ {
			bootStrapAddr, IPerr = net.LookupIP(bootstrapNode)
			if IPerr == nil && len(bootStrapAddr) > 0 {
				break
			}
			fmt.Printf("LookupIP error (attempt %d): %v\n", i+1, IPerr)
			time.Sleep(2 * time.Second)
		}
		if IPerr != nil || len(bootStrapAddr) == 0 {
			fmt.Fprintf(os.Stderr, "Failed to resolve bootstrap node after retries: %v\n", IPerr)
			os.Exit(1)
		}
		bootstrapIP = bootStrapAddr[0].String() + ":" + "9001"

		k, kadErr = kademlia.InitKademlia(port, false, bootstrapIP)
		if kadErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize Kademlia: %v\n", kadErr)
			os.Exit(1)
		}
	}

	k.Server.RunServer()
	if isBootstrap == "TRUE" {
		time.Sleep(15 * time.Second)
		ans, err := k.Client.SendStoreMessage([]byte("test123"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to SendStoreMessage %v\n", err)
		} else {
			fmt.Println("Sent SendStoreMessage: ", ans)
		}
	}

	if isBootstrap == "FALSE" {
		time.Sleep(5 * time.Second)
		bootstrapContact := kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000000"), bootstrapIP)
		_, err := k.Client.SendPingMessage(bootstrapContact)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to ping bootstrap node: %v\n", err)
		} else {
			fmt.Println("Sent ping to: ", bootstrapContact)
		}

		// Do iterative find node
		time.Sleep(10 * time.Second)
		k.Node.IterativeFindNode(k.Node.GetSelfContact().ID)

		time.Sleep(30 * time.Second)
		k.Node.PrintRoutingTable()
	}

	if os.Getenv("ENABLECLI") == "TRUE" {
		k.Node.Cli()
	}
	select {} // keep running
}
