package kademlia

import "net"

type UDPNetwork struct {
	conn *net.UDPConn
}

type Network interface {
	SendMessage(addr string, data []byte) error
	ReceiveMessage() (addr string, data []byte, err error)
	Close() error
	GetConn() string
}

func NewUDPNetwork(localAddr string) (*UDPNetwork, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	return &UDPNetwork{conn: conn}, nil
}

func (u *UDPNetwork) GetConn() string {
	return u.conn.LocalAddr().String()
}

func (u *UDPNetwork) ReceiveMessage() (string, []byte, error) {
	buf := make([]byte, 8192)
	n, addr, err := u.conn.ReadFromUDP(buf)
	if err != nil {
		return "", nil, err
	}
	return addr.String(), buf[:n], nil
}

func (u *UDPNetwork) SendMessage(addr string, data []byte) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	_, err = u.conn.WriteToUDP(data, udpAddr)
	return err
}

func (u *UDPNetwork) Close() error {
	return u.conn.Close()
}
