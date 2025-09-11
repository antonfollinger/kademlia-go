package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
)

const (
	IncomingBufferSize int = 256
	OutgoingBufferSize int = 64
)

type Server struct {
	conn     *net.UDPConn
	incoming chan RPCMessage
	outgoing chan RPCMessage
}

func InitServer(ip string) (*Server, error) {
	fmt.Println("addr: ", ip)
	udpAddr, err := net.ResolveUDPAddr("udp", ip)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		conn:     conn,
		incoming: make(chan RPCMessage, IncomingBufferSize),
		outgoing: make(chan RPCMessage, OutgoingBufferSize),
	}

	return s, nil
}

func (s *Server) RunServer() {
	fmt.Println("RunServer()")
	go s.listen()
	go s.handleIncoming()
	go s.respond()
}

func (s *Server) listen() {
	buf := make([]byte, 4096)
	for {
		n, addr, err := s.conn.ReadFromUDP(buf)
		fmt.Println("Found RPC")
		if err != nil {
			fmt.Println("listen error:", err)
			continue
		}
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
			fmt.Println("handlePing()")
			resp = s.handlePing(rpc)
		default:
			resp = *CreateRPCMessage("ERROR", Payload{})
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
	resp := CreateRPCMessage("OK", Payload{})

	// Ensure same PID
	PID := rpc.PacketID
	resp.PacketID = PID

	return *resp
}
