package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
)

type IncomingRPC struct {
	RPC  RPCMessage
	Addr *net.UDPAddr
}

type OutgoingRPC struct {
	RPC  RPCMessage
	Addr *net.UDPAddr
}

const (
	IncomingBufferSize int = 4096
	OutgoingBufferSize int = 1024
)

type Server struct {
	node     NodeAPI
	conn     *net.UDPConn
	incoming chan IncomingRPC
	outgoing chan OutgoingRPC
}

func InitServer(node NodeAPI) (*Server, error) {
	ip := node.GetSelfContact().Address
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
		incoming: make(chan IncomingRPC, IncomingBufferSize),
		outgoing: make(chan OutgoingRPC, OutgoingBufferSize),
	}
	fmt.Println("Server listening on: ", ip)
	return s, nil
}

func (s *Server) RunServer() {
	go s.listen()
	go s.handleIncoming()
	go s.respond()
}

func (s *Server) listen() {
	buf := make([]byte, 4096)
	for {
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("listen error:", err)
			continue
		}
		fmt.Printf("Server found RPC from %v, bytes read: %d\n\n", addr, n)
		var rpc RPCMessage

		if err := json.Unmarshal(buf[:n], &rpc); err != nil {
			fmt.Println("Unmarshal error:", err)
			continue
		}
		if rpc.Query {
			fmt.Printf("RPC INFO: %+v\n\n", rpc)
			s.incoming <- IncomingRPC{RPC: rpc, Addr: addr}
		}
	}
}

func (s *Server) handleIncoming() {
	for in := range s.incoming {
		var resp RPCMessage
		switch in.RPC.Type {
		case "PING":
			resp = s.handlePing(in.RPC)
		default:
			resp = *NewRPCMessage("ERROR", Payload{SourceContact: in.RPC.Payload.SourceContact}, false)
		}
		s.outgoing <- OutgoingRPC{RPC: resp, Addr: in.Addr}
	}
}

func (s *Server) respond() {
	for {
		out := <-s.outgoing
		fmt.Println("Responding to:", out.Addr)
		data, _ := json.Marshal(out.RPC)
		_, _ = s.conn.WriteToUDP(data, out.Addr)
	}
}

func (s *Server) handlePing(rpc RPCMessage) RPCMessage {
	resp := NewRPCMessage("PONG", Payload{SourceContact: rpc.Payload.SourceContact}, false)

	// Ensure same PID
	PID := rpc.PacketID
	resp.PacketID = PID

	return *resp
}
