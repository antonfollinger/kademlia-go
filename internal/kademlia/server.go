package kademlia

import (
	"net"
)

const (
	IncomingBufferSize int = 256
	OutgoingBufferSize int = 64
)

type Message struct {
	SourceAddr      string
	DestinationAddr string
	Data            string
}

type Server struct {
	conn     *net.UDPConn
	incoming chan Message
	outgoing chan Message
}

func InitServer(ip string) (*Server, error) {
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
		incoming: make(chan Message, IncomingBufferSize),
		outgoing: make(chan Message, OutgoingBufferSize),
	}

	return s, nil
}

func (s *Server) Run() {
	go s.listen()
	go s.respond()
}

func (s *Server) listen() {
	buf := make([]byte, 1024)
	for {
		select {
		default:
			n, _, err := s.conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}
			msg := Message{}
			s.incoming <- msg
		}
	}
}

func (s *Server) respond() {

}
