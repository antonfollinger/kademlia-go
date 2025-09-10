package kademlia

import "net"

type Kademlia struct {
	addr        string
	node        *Node
	client      *Client
	server      *Server
	isBootstrap bool
}

func InitKademlia(bootStrap bool) (*Kademlia, error) {
	k := &Kademlia{}
	ip := k.getLocalIP()

	k.addr = ip
	k.isBootstrap = bootStrap

	// Node
	var nodeErr error
	k.node, nodeErr = InitNode(bootStrap, ip)
	if nodeErr != nil {
		return nil, nodeErr
	}

	// Client
	var clientErr error
	k.client, clientErr = InitClient(ip)
	if clientErr != nil {
		return nil, clientErr
	}

	// Server
	var serverErr error
	k.server, serverErr = InitServer(ip)
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
