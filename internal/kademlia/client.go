package kademlia

const (
	ClientBufferSize int = 64
)

type Client struct {
	node     *Kademlia
	addr     string
	request  chan string
	response chan string
}

func InitClient(node *Kademlia, ip string) (*Client, error) {
	c := &Client{
		node:     node,
		request:  make(chan string, ClientBufferSize),
		response: make(chan string, ClientBufferSize),
	}
	c.addr = ip

	return c, nil
}

/*
func (c *Client) SendPingMessage(ip string) error {
	addr, err := net.ResolveUDPAddr("udp", ip)
	fmt.Println("addr: ", addr.String())
	if err != nil {
		fmt.Println("Resolve error: ", err)
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	fmt.Println("conn: ", conn.LocalAddr())
	if err != nil {
		fmt.Println("Dial error: ", err)
		return err
	}
	defer conn.Close()

	msg := CreateRPCMessage("PING", Payload{SourceContact: &c.node.RoutingTable.me})
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Marshal error: ", err)
		return err
	}

	_, err = conn.Write(data)
	fmt.Println("PING sent to: ", ip)
	if err != nil {
		return err
	}

	return nil
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
*/
