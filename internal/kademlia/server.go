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
	IncomingBufferSize int = 8192
	OutgoingBufferSize int = 8192
)

type Server struct {
	node     NodeAPI
	network  Network
	incoming chan IncomingRPC
	outgoing chan OutgoingRPC
}

func InitServer(node NodeAPI, network Network) (*Server, error) {
	s := &Server{
		node:     node,
		network:  network,
		incoming: make(chan IncomingRPC, IncomingBufferSize),
		outgoing: make(chan OutgoingRPC, OutgoingBufferSize),
	}
	fmt.Println("Server listening on: ", network.GetConn())
	return s, nil
}

func (s *Server) RunServer() {
	go s.listen()
	go s.handleIncoming()
	go s.respond()
}

func (s *Server) listen() {
	for {
		addrStr, data, err := s.network.ReceiveMessage()
		if err != nil {
			fmt.Println("listen error:", err)
			continue
		}

		var rpc RPCMessage
		if err := json.Unmarshal(data, &rpc); err != nil {
			fmt.Println("Unmarshal error:", err)
			continue
		}

		var addr *net.UDPAddr
		if addrStr != "" {
			addr, _ = net.ResolveUDPAddr("udp", addrStr)
		}

		if rpc.Query {
			select {
			case s.incoming <- IncomingRPC{RPC: rpc, Addr: addr}:
				// Packet accepted
			default:
				fmt.Printf("DROPPED PACKET from %v: incoming channel overflow\n", addr)
			}
		}
	}
}

func (s *Server) handleIncoming() {
	for in := range s.incoming {
		go s.processRequest(in)
	}
}

func (s *Server) processRequest(in IncomingRPC) {
	var resp RPCMessage
	switch in.RPC.Type {
	case "PING":
		resp = *NewRPCMessage("PONG", Payload{
			TargetContact: in.RPC.Payload.SourceContact,
		}, false)
	case "FIND_NODE":
		target := NewKademliaID(in.RPC.Payload.Key)
		contacts := s.node.LookupClosestContacts(NewContact(target, ""))
		resp = *NewRPCMessage("FIND_NODE", Payload{
			Contacts:      contacts,
			TargetContact: in.RPC.Payload.SourceContact,
		}, false)
	case "STORE":
		s.node.Store(in.RPC.Payload.Key, in.RPC.Payload.Data)
		contacts := s.node.GetSelfContact()
		resp = *NewRPCMessage("STORE", Payload{
			Contacts:      []Contact{contacts},
			TargetContact: in.RPC.Payload.SourceContact,
			Key:           in.RPC.Payload.Key,
		}, false)
	case "FIND_VALUE":
		value := s.node.LookupData(in.RPC.Payload.Key)
		resp = *NewRPCMessage("FIND_VALUE", Payload{
			Data:          value,
			TargetContact: in.RPC.Payload.SourceContact,
		}, false)
	default:
		resp = *NewRPCMessage("ERROR", Payload{TargetContact: in.RPC.Payload.TargetContact}, false)
	}

	// Add the requesting node
	s.node.AddContact(in.RPC.Payload.SourceContact)

	// Ensure same PID
	PID := in.RPC.PacketID
	resp.PacketID = PID

	resp.Payload.SourceContact = s.node.GetSelfContact()

	s.outgoing <- OutgoingRPC{RPC: resp, Addr: in.Addr}
}

func (s *Server) respond() {
	for {
		out := <-s.outgoing
		data, _ := json.Marshal(out.RPC)
		if out.Addr != nil {
			_ = s.network.SendMessage(out.Addr.String(), data)
		}
	}
}
