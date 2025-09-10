package kademlia

import (
	"fmt"
	"net"
)

const (
	ClientBufferSize int = 64
)

type Client struct {
	addr     string
	request  chan string
	response chan string
}

func InitClient(ip string) (*Client, error) {
	c := &Client{
		request:  make(chan string, ClientBufferSize),
		response: make(chan string, ClientBufferSize),
	}
	c.addr = ip

	return c, nil
}

func (server *Server) SendPingMessage(contact *Contact) {
	server.sendMessage(contact, "PING")
}

func (server *Server) SendFindContactMessage(contact *Contact) {
	server.sendMessage(contact, "FIND_CONTACT")
}

func (server *Server) SendFindDataMessage(hash string) {
	// For demo, just print
	fmt.Println("SendFindDataMessage called with hash:", hash)
}

func (server *Server) SendStoreMessage(data []byte) {
	// For demo, just print
	fmt.Println("SendStoreMessage called with data:", string(data))
}

func (server *Server) sendMessage(contact *Contact, msg string) {
	addr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		fmt.Println("ResolveUDPAddr error:", err)
		return
	}
	_, err = server.conn.WriteToUDP([]byte(msg), addr)
	if err != nil {
		fmt.Println("WriteToUDP error:", err)
	}
}
