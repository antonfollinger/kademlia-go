package kademlia

import (
	"encoding/json"
	"fmt"
)

type IncomingRPC struct {
	RPC  RPCMessage
	Addr string
}

type OutgoingRPC struct {
	RPC  RPCMessage
	Addr string
}

const (
	IncomingBufferSize int = 256
	OutgoingBufferSize int = 128
	workerCount        int = 5
)

type Server struct {
	node     NodeAPI
	network  Network
	incoming chan IncomingRPC
	outgoing chan OutgoingRPC
	done     chan struct{}
}

func InitServer(node NodeAPI, network Network) (*Server, error) {
	s := &Server{
		node:     node,
		network:  network,
		incoming: make(chan IncomingRPC, IncomingBufferSize),
		outgoing: make(chan OutgoingRPC, OutgoingBufferSize),
		done:     make(chan struct{}),
	}

	s.RunServer()

	return s, nil
}

func (s *Server) RunServer() {
	go s.listen()
	go s.respond()
	for range workerCount {
		go func() {
			for {
				select {
				case <-s.done:
					return
				case in := <-s.incoming:
					s.processRequest(in)
				}
			}
		}()
	}
}

func (s *Server) listen() {
	for {
		select {
		case <-s.done:
			return
		default:
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

			if rpc.Query {
				select {
				case s.incoming <- IncomingRPC{RPC: rpc, Addr: addrStr}:
					// Packet accepted
				default:
					fmt.Printf("DROPPED PACKET from %v: incoming channel overflow\n", addrStr)
				}
			}
		}
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
		select {
		case <-s.done:
			return
		case out := <-s.outgoing:
			data, _ := json.Marshal(out.RPC)
			if out.Addr != "" {
				_ = s.network.SendMessage(out.Addr, data)
			}
		}
	}
}

// Close gracefully shuts down the server and its network connection
func (s *Server) Close() error {
	close(s.done)
	close(s.incoming)
	close(s.outgoing)
	return s.network.Close()
}
