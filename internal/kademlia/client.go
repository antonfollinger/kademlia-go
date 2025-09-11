package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
)

const (
	ClientBufferSize int = 64
)

type Client struct {
	conn     *net.UDPConn
	addr     string
	request  chan string
	response chan string
}

func InitClient(ip string) (*Client, error) {
	addr, _ := net.ResolveUDPAddr("udp", ":0") // :0 = pick random available port
	conn, _ := net.ListenUDP("udp", addr)
	c := &Client{
		conn:     conn,
		request:  make(chan string, ClientBufferSize),
		response: make(chan string, ClientBufferSize),
	}
	c.addr = ip

	return c, nil
}

func (c *Client) SendPingMessage(ip string) {
	fmt.Println("Sending ping")
	// Create the RPC packet
	packet := CreateRPCMessage("PING", Payload{})

	data, err := json.Marshal(packet)
	if err != nil {
		return
	}

	// Resolve target address
	addr, err := net.ResolveUDPAddr("udp", ip)
	if err != nil {
		return
	}

	// Send via existing UDPConn
	_, _ = c.conn.WriteToUDP(data, addr)
}

func (server *Server) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (server *Server) SendFindDataMessage(hash string) {
	// TODO
}

func (server *Server) SendStoreMessage(data []byte) {
	// TODO
}
