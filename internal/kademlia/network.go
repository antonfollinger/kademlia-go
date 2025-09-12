package kademlia

import (
	"net"
)

type Network struct {
	Kademlia *Kademlia
	Addr     string
	Server   *Server
	Client   *Client
}

func initNetwork(k *Kademlia) (*Network, error) {

	n := &Network{Kademlia: k}

	ip := n.GetLocalIP()

	// Client
	var clientErr error
	n.Client, clientErr = InitClient(k, ip)
	if clientErr != nil {
		return nil, clientErr
	}

	// Server
	var serverErr error
	n.Server, serverErr = InitServer(k, ip)
	if serverErr != nil {
		return nil, serverErr
	}

	return n, nil
}

func (network *Network) GetLocalIP() string {
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

/*
func Listen(ip string, port int) (*Network, error) {
	addr, err := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	network := &Network{Conn: conn}
	go network.listenLoop()
	return network, nil
}

func (network *Network) listenLoop() {
	buf := make([]byte, 2048) // bigger buffer for JSON
	for {
		n, addr, err := network.Conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading UDP:", err)
			continue
		}

		// Decode into RPCMessage
		var msg RPCMessage
		if err := json.Unmarshal(buf[:n], &msg); err != nil {
			fmt.Println("Invalid JSON from", addr, ":", string(buf[:n]))
			continue
		}

		// Hand off to the RPC handler
		go network.handleRPC(&msg)
	}
}

func (network *Network) SendMessage(contact *Contact, msg *RPCMessage) error {
	// Marshal RPCMessage into JSON
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal RPCMessage: %w", err)
	}

	// Build UDP address from Contact
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", contact.Address, contact.Port))
	if err != nil {
		return fmt.Errorf("failed to resolve UDP addr: %w", err)
	}

	// Send JSON bytes
	_, err = network.Conn.WriteToUDP(data, addr)
	if err != nil {
		return fmt.Errorf("failed to send UDP message: %w", err)
	}

	return nil
}
*/
