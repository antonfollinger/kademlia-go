package kademlia

import (
	"fmt"
	"net"
	"strconv"
)

type Contact struct {
	IP   string
	Port int
}

type Network struct {
	Conn *net.UDPConn
}

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
	buf := make([]byte, 1024)
	for {
		n, addr, err := network.Conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}
		msg := string(buf[:n])
		fmt.Printf("Received from %s: %s\n", addr, msg)
		if msg == "PING" {
			// Send PONG back to sender
			_, err := network.Conn.WriteToUDP([]byte("PONG"), addr)
			if err != nil {
				fmt.Println("Error sending PONG:", err)
			} else {
				fmt.Printf("Sent PONG to %s\n", addr)
			}
		}
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	network.sendMessage(contact, "PING")
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	network.sendMessage(contact, "FIND_CONTACT")
}

func (network *Network) SendFindDataMessage(hash string) {
	// For demo, just print
	fmt.Println("SendFindDataMessage called with hash:", hash)
}

func (network *Network) SendStoreMessage(data []byte) {
	// For demo, just print
	fmt.Println("SendStoreMessage called with data:", string(data))
}

func (network *Network) sendMessage(contact *Contact, msg string) {
	addr, err := net.ResolveUDPAddr("udp", contact.IP+":"+strconv.Itoa(contact.Port))
	if err != nil {
		fmt.Println("ResolveUDPAddr error:", err)
		return
	}
	_, err = network.Conn.WriteToUDP([]byte(msg), addr)
	if err != nil {
		fmt.Println("WriteToUDP error:", err)
	}
}
