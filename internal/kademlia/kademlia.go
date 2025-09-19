package kademlia

import (
	"fmt"
	"net"
	"os"
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

	bootstrapID := NewKademliaID("0000000000000000000000000000000000000000")

	if !k.Node.RoutingTable.me.ID.Equals(bootstrapID) {

		_, err1 := k.Client.SendPingMessage(k.Node.RoutingTable.FindClosestContacts(bootstrapID, 1)[0])
		if err1 != nil {
			fmt.Printf("JoinNetwork: Unable to reach bootstrap node: %v\n", err1)

			// Retry pinging bootstrap a few times
			for i := 0; i < 4; i++ {
				_, err1 = k.Client.SendPingMessage(k.Node.RoutingTable.FindClosestContacts(bootstrapID, 1)[0])
				if err1 == nil {
					break
				}
			}
			if err1 != nil {
				return k, err1
			}
		}

	}
	// Integrate JoinNetwork for both bootstrap and peer nodes
	err := k.Node.JoinNetwork()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to join network: %v\n", err)
	}

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
