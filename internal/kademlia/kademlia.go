package kademlia

import (
	"fmt"
	"net"
)

type Kademlia struct {
	Node   *Node
	Server *Server
	Client *Client
}

func InitKademlia(port string, bootstrap bool, bootstrapIP string) (*Kademlia, error) {

	k := &Kademlia{}
	ip := GetLocalIP() + ":" + port

	fmt.Println("Local_ip: ", ip)
	fmt.Println("Bootstrap IP: ", bootstrapIP)

	// Node
	var nodeErr error
	k.Node, nodeErr = InitNode(bootstrap, ip, bootstrapIP)
	if nodeErr != nil {
		return nil, nodeErr
	}

	// Client
	var clientErr error
	k.Client, clientErr = InitClient(k.Node)
	if clientErr != nil {
		return nil, clientErr
	}

	// Server
	var serverErr error
	k.Server, serverErr = InitServer(k.Node)
	if serverErr != nil {
		return nil, serverErr
	}

	k.Node.SetClient(k.Client)

	return k, nil
}

func GetLocalIP() string {
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
