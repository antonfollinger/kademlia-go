package kademlia

import (
	"log"
	"math/rand"
	"net"
	"time"
)

type KademliaConfig struct {
	SkipBootstrapPing    bool
	BootstrapPingRetries int
	BootstrapPingDelayMs int
	isMockNetwork        bool
	MockNetworkRegistry  *MockRegistry
}

type KademliaOption func(*KademliaConfig)

func WithSkipBootstrapPing(skip bool) KademliaOption {
	return func(cfg *KademliaConfig) {
		cfg.SkipBootstrapPing = skip
	}
}

type Kademlia struct {
	Node   *Node
	Server *Server
	Client *Client
}

func InitKademlia(port string, bootstrap bool, bootstrapIP string, opts ...KademliaOption) (*Kademlia, error) {
	cfg := &KademliaConfig{
		SkipBootstrapPing:    false,
		BootstrapPingRetries: 5,
		BootstrapPingDelayMs: 500,
		isMockNetwork:        false,
		MockNetworkRegistry:  nil,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	k := &Kademlia{}
	var ip string
	if cfg.isMockNetwork {
		ip = "127.0.0.1" + ":" + port
	} else {
		ip = GetLocalIP() + ":" + port
	}

	log.Println("Local_ip: ", ip)
	log.Println("Bootstrap IP: ", bootstrapIP)

	// Node
	var nodeErr error
	k.Node, nodeErr = InitNode(bootstrap, ip, bootstrapIP)
	if nodeErr != nil {
		return nil, nodeErr
	}

	// Network selection
	var clientNet Network
	var serverNet Network
	if cfg.isMockNetwork {
		// Use unique addresses for client and server to avoid channel sharing
		clientAddr := "127.0.0.1" + ":" + port + ":client"
		clientNet = NewMockNetwork(clientAddr, cfg.MockNetworkRegistry)
		serverNet = NewMockNetwork(ip, cfg.MockNetworkRegistry)
	} else {
		var err error
		clientNet, err = NewUDPNetwork("") // ephemeral port
		if err != nil {
			return nil, err
		}
		serverNet, err = NewUDPNetwork(ip)
		if err != nil {
			return nil, err
		}
	}

	// Client
	var clientErr error
	k.Client, clientErr = InitClient(k.Node, clientNet)
	if clientErr != nil {
		return nil, clientErr
	}

	// Server
	var serverErr error
	k.Server, serverErr = InitServer(k.Node, serverNet)
	if serverErr != nil {
		return nil, serverErr
	}

	k.Node.SetClient(k.Client)

	bootstrapID := NewKademliaID("0000000000000000000000000000000000000000")

	if !k.Node.RoutingTable.me.ID.Equals(bootstrapID) && !cfg.SkipBootstrapPing {
		// Random delay to reduce package drops
		time.Sleep(time.Duration(rand.Intn(cfg.BootstrapPingDelayMs)) * time.Millisecond)

		// Retry pinging bootstrap a few times
		var err1 error
		for i := 0; i < cfg.BootstrapPingRetries; i++ {
			c := k.Node.RoutingTable.FindClosestContacts(bootstrapID, 1)[0]
			_, err1 = k.Client.SendPingMessage(c)
			if err1 == nil {
				break
			}
			time.Sleep(time.Duration(rand.Intn(cfg.BootstrapPingDelayMs)) * time.Millisecond)
		}
		if err1 != nil {
			println(ip, " failed to ping bootstrap node on", bootstrapIP)
		}
	}

	if !cfg.SkipBootstrapPing {
		// Integrate JoinNetwork for both bootstrap and peer nodes
		err := k.Node.JoinNetwork()
		if err != nil {
			log.Printf("failed to join network: %v\n", err)
		}
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
