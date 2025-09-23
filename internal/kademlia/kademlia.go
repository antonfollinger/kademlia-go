package kademlia

import (
	"log"
	"net"
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
