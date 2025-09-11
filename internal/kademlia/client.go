package kademlia

import (
	"encoding/json"
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

func (server *Server) SendPingMessage(ip string) {
	conn, err := net.Dial("udp", ip)
	if err != nil {
		// handle error, e.g., log or return
		return
	}
	defer conn.Close()

	packet := CreateRPCMessage("PING", Payload{})
	data, err := json.Marshal(packet)
	if err != nil {
		// handle error
		return
	}

	_, _ = conn.Write(data)
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
