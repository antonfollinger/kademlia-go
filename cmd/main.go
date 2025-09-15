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
		// Bootstrap Node info
		bootstrapNode := os.Getenv("BOOTSTRAPNODE")
		bootStrapAddr, IPerr := net.LookupIP(bootstrapNode)
		if IPerr != nil {
			fmt.Println("LookupIP error", IPerr)
		}
		bootstrapIP = bootStrapAddr[0].String() + ":" + "9001"

		k, kadErr = kademlia.InitKademlia(port, false, bootstrapIP)
		if kadErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize Kademlia: %v\n", kadErr)
			os.Exit(1)
		}
	}

	k.Server.RunServer()

	if isBootstrap == "FALSE" {
		time.Sleep(2 * time.Second)
		bootstrapContact := kademlia.NewContact(kademlia.NewKademliaID("0000000000000000000000000000000000000000"), bootstrapIP)
		_, err := k.Client.SendPingMessage(bootstrapContact)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to ping bootstrap node: %v\n", err)
		} else {
			fmt.Println("Sent ping to: ", bootstrapContact)
		}
	}

	select {} // keep running
}
