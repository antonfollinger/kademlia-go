package kademlia

import (
	"fmt"
	"net"
)

type Kademlia struct {
	addr        string
	Node        *Node
	Client      *Client
	Server      *Server
	isBootstrap bool
}

func InitKademlia(bootStrap bool, port string) (*Kademlia, error) {
	k := &Kademlia{}

	var ip string
	if bootStrap {
		ip = "0.0.0.0:" + port
	} else {
		ip = k.getLocalIP() + ":" + port
	}
	k.addr = ip
	k.isBootstrap = bootStrap

	// Node
	var nodeErr error
	k.Node, nodeErr = InitNode(bootStrap, ip)
	if nodeErr != nil {
		return nil, nodeErr
	}

	fmt.Println("ME: ", k.Node.RoutingTable.Me.String())

	// Client
	var clientErr error
	k.Client, clientErr = InitClient(k.Node, ip)
	if clientErr != nil {
		return nil, clientErr
	}

	// Server
	var serverErr error
	k.Server, serverErr = InitServer(k.Node, ip)
	if serverErr != nil {
		return nil, serverErr
	}

	return k, nil
}

func (kademlia *Kademlia) getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
