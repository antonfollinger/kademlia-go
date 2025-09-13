package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
)

const (
	IncomingBufferSize int = 4096
	OutgoingBufferSize int = 1024
)

type Server struct {
	node     NodeAPI
	conn     *net.UDPConn
	incoming chan RPCMessage
	outgoing chan RPCMessage
}

func InitServer(node *Kademlia, ip string) (*Server, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", ip)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		node:     node,
		conn:     conn,
		incoming: make(chan RPCMessage, IncomingBufferSize),
		outgoing: make(chan RPCMessage, OutgoingBufferSize),
	}
	fmt.Println("Listening on: ", ip)
	return s, nil
}

func (s *Server) RunServer() {
	fmt.Println("Server running...")
	go s.listen()
	go s.handleIncoming()
	go s.respond()
}

func (s *Server) listen() {
	fmt.Println("Listening...")
	buf := make([]byte, 4096)
	for {
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("listen error:", err)
			continue
		}
		fmt.Printf("Found RPC from %v, bytes read: %d\n", addr, n)
		var rpc RPCMessage

		if err := json.Unmarshal(buf[:n], &rpc); err != nil {
			continue
		}

		if rpc.Payload.SourceContact.Address == "" {
			rpc.Payload.SourceContact.Address = addr.String()
		}

		s.incoming <- rpc
	}
}

func (s *Server) handleIncoming() {
	for rpc := range s.incoming {
		var resp RPCMessage
		switch rpc.Type {
		case "PING":
			resp = s.handlePing(rpc)
		default:
			resp = *NewRPCMessage("ERROR", Payload{}, false)
		}
		s.outgoing <- resp
	}
}

func (s *Server) respond() {
	for {
		rpc := <-s.outgoing
		target := rpc.Payload.SourceContact.Address // must carry destination
		data, _ := json.Marshal(rpc)
		addr, _ := net.ResolveUDPAddr("udp", target)
		_, _ = s.conn.WriteToUDP(data, addr)
	}
}

func (s *Server) handlePing(rpc RPCMessage) RPCMessage {
	resp := NewRPCMessage("OK", Payload{SourceContact: rpc.Payload.SourceContact}, false)

	// Ensure same PID
	PID := rpc.PacketID
	resp.PacketID = PID

	return *resp
}
