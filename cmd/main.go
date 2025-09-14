package main

import (
	"fmt"
	"net"
	"os"

	"github.com/antonfollinger/kademlia_go/internal/kademlia"
)

func main() {
	isBootstrap := os.Getenv("ISBOOTSTRAP")
	port := os.Getenv("PORT")

	var k *kademlia.Kademlia
	var kadErr error

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
		bootstrapIP := bootStrapAddr[0].String() + ":" + "9001"

		k, kadErr = kademlia.InitKademlia(port, false, bootstrapIP)
		if kadErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize Kademlia: %v\n", kadErr)
			os.Exit(1)
		}
	}

	k.Server.RunServer()

	select {} // keep running
}
