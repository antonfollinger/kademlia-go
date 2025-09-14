package kademlia

const (
	ClientBufferSize int = 64
)

type Client struct {
	node     NodeAPI
	request  chan string
	response chan string
}

func InitClient(node NodeAPI) (*Client, error) {
	c := &Client{
		node:     node,
		request:  make(chan string, ClientBufferSize),
		response: make(chan string, ClientBufferSize),
	}

	return c, nil
}

/*
func (network *Network) SendPingMessage(msg *RPCMessage) error {
	// Build payload with my own contact
	payload := Payload{
		SourceContact: &network.Kademlia.RoutingTable.me,
	}

	if msg.Query {
		if msg.Payload.SourceContact != nil {
			network.Kademlia.RoutingTable.AddContact(*msg.Payload.SourceContact)
			reply := NewRPCMessage("PING", payload, false)
			reply.PacketID = msg.PacketID // match request ID
			fmt.Printf("Got PING from %s:%d (ID=%s, PacketID=%s)\n",
				msg.Payload.SourceContact.Address,
				msg.Payload.SourceContact.Port,
				msg.Payload.SourceContact.ID.String(),
				msg.PacketID)
			_ = network.SendMessage(msg.Payload.SourceContact, reply)
		}
	} else {
		if msg.Payload.SourceContact != nil {
			network.Kademlia.RoutingTable.AddContact(*msg.Payload.SourceContact)
			fmt.Printf("Got PONG from %s:%d (ID=%s, PacketID=%s)\n",
				msg.Payload.SourceContact.Address,
				msg.Payload.SourceContact.Port,
				msg.Payload.SourceContact.ID.String(),
				msg.PacketID)
		}
	}

	return nil
}

func (network *Network) SendFindContactMessage(msg *RPCMessage) {
	if msg.Query {
		contacts := network.Kademlia.RoutingTable.FindClosestContacts(msg.Payload.TargetContact.ID, 8)

		payload := Payload{
			Contacts:      contacts,
			SourceContact: &network.Kademlia.RoutingTable.me,
		}
		reply := NewRPCMessage("FIND_NODE", payload, false)
		reply.PacketID = msg.PacketID // preserve request ID

		if msg.Payload.SourceContact != nil {
			_ = network.SendMessage(msg.Payload.SourceContact, reply)
		}
	} else {
		if msg.Payload.SourceContact != nil {
			network.Kademlia.RoutingTable.AddContact(*msg.Payload.SourceContact)
		}

		if msg.Payload.Contacts != nil {
			fmt.Printf("Got FIND_NODE response with %d contacts (PacketID=%s)\n",
				len(msg.Payload.Contacts), msg.PacketID)
			for _, contact := range msg.Payload.Contacts {
				network.Kademlia.RoutingTable.AddContact(contact)
			}
		}
	}
}

func (network *Network) SendStoreMessage(msg *RPCMessage) {
	if msg.Query {
		network.Kademlia.Store(msg.Payload.Key, msg.Payload.Data)

		payload := Payload{
			SourceContact: &network.Kademlia.RoutingTable.me,
		}
		msgout := NewRPCMessage("STORE", payload, false)
		network.SendMessage(msg.Payload.SourceContact, msgout)

	} else {
		if msg.Payload.SourceContact != nil {
			network.Kademlia.RoutingTable.AddContact(*msg.Payload.SourceContact)

			network.Kademlia.RoutingTable.AddContact(*msg.Payload.SourceContact)
			fmt.Printf("Got STORE ACK from %s:%d (ID=%s, PacketID=%s)\n",
				msg.Payload.SourceContact.Address,
				msg.Payload.SourceContact.Port,
				msg.Payload.SourceContact.ID.String(),
				msg.PacketID)
			network.Kademlia.RoutingTable.AddContact(*msg.Payload.SourceContact)
		}
	}

}

func (network *Network) SendFindValueMessage(msg *RPCMessage) {
	if msg.Query {
		payload := Payload{
			Data:          network.Kademlia.LookupData(msg.Payload.Key),
			SourceContact: &network.Kademlia.RoutingTable.me,
		}
		msgout := NewRPCMessage("FIND_VALUE", payload, false)
		network.SendMessage(msg.Payload.SourceContact, msgout)
	} else {
		network.Kademlia.RoutingTable.AddContact(*msg.Payload.SourceContact)

		if msg.Payload.Data != nil {
			fmt.Printf("Got FIND_VALUE response with data: %s (PacketID=%s)\n",
				string(msg.Payload.Data),
				msg.PacketID)
		}
	}
}

func (network *Network) handleRPC(msg *RPCMessage) {
	switch msg.Type {
	case "PING":
		network.SendPingMessage(msg)
	case "FIND_NODE":
		network.SendFindContactMessage(msg)
	case "STORE":
		network.SendStoreMessage(msg)
	case "FIND_VALUE":
		network.SendFindValueMessage(msg)
	default:
		fmt.Println("Unknown RPC message type:", msg.Type)
	}
}
*/

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
*/
