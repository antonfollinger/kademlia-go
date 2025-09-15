package kademlia

import "fmt"

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

func (client *Client) SendPingMessage(msg *RPCMessage) error {
	// Build payload with my own contact
	payload := Payload{
		SourceContact: client.node.GetSelfContact(),
	}

	if msg.Query {
		if msg.Payload.SourceContact != nil {
			client.node.AddContact(msg.Payload.SourceContact)
			reply := NewRPCMessage("PING", payload, false)
			reply.PacketID = msg.PacketID // match request ID
			fmt.Printf("Got PING from %s:%d (ID=%s, PacketID=%s)\n",
				msg.Payload.SourceContact.Address,
				msg.Payload.SourceContact.ID.String(),
				msg.PacketID)
			_ = client.SendMessage(msg.Payload.SourceContact, reply)
		}
	} else {
		if msg.Payload.SourceContact != nil {
			client.node.AddContact(msg.Payload.SourceContact)
			fmt.Printf("Got PONG from %s:%d (ID=%s, PacketID=%s)\n",
				msg.Payload.SourceContact.Address,
				msg.Payload.SourceContact.ID.String(),
				msg.PacketID)
		}
	}

	return nil
}

func (client *Client) SendFindContactMessage(msg *RPCMessage) {
	if msg.Query {
		contacts := client.node.LookupContact(msg.Payload.TargetContact)

		payload := Payload{
			Contacts:      contacts,
			SourceContact: client.node.GetSelfContact(),
		}
		reply := NewRPCMessage("FIND_NODE", payload, false)
		reply.PacketID = msg.PacketID // preserve request ID

		if msg.Payload.SourceContact != nil {
			_ = client.SendMessage(msg.Payload.SourceContact, reply)
		}
	} else {
		if msg.Payload.SourceContact != nil {
			client.node.AddContact(msg.Payload.SourceContact)
		}

		if msg.Payload.Contacts != nil {
			fmt.Printf("Got FIND_NODE response with %d contacts (PacketID=%s)\n",
				len(msg.Payload.Contacts), msg.PacketID)
			for _, contact := range msg.Payload.Contacts {
				client.node.AddContact(contact)
			}
		}
	}
}

func (client *Client) SendStoreMessage(msg *RPCMessage) {
	if msg.Query {
		client.node.Store(msg.Payload.Key, msg.Payload.Data)

		payload := Payload{
			SourceContact: client.node.GetSelfContact(),
		}
		msgout := NewRPCMessage("STORE", payload, false)
		client.SendMessage(msg.Payload.SourceContact, msgout)

	} else {
		if msg.Payload.SourceContact != nil {
			client.node.AddContact(msg.Payload.SourceContact)

			client.node.AddContact(msg.Payload.SourceContact)
			fmt.Printf("Got STORE ACK from %s:%d (ID=%s, PacketID=%s)\n",
				msg.Payload.SourceContact.Address,
				msg.Payload.SourceContact.ID.String(),
				msg.PacketID)
			client.node.AddContact(msg.Payload.SourceContact)
		}
	}

}

func (client *Client) SendFindValueMessage(msg *RPCMessage) {
	if msg.Query {
		payload := Payload{
			Data:          client.node.LookupData(msg.Payload.Key),
			SourceContact: client.node.GetSelfContact(),
		}
		msgout := NewRPCMessage("FIND_VALUE", payload, false)
		client.SendMessage(msg.Payload.SourceContact, msgout)
	} else {
		client.node.AddContact(msg.Payload.SourceContact)

		if msg.Payload.Data != nil {
			fmt.Printf("Got FIND_VALUE response with data: %s (PacketID=%s)\n",
				string(msg.Payload.Data),
				msg.PacketID)
		}
	}
}

func (client *Client) handleRPC(msg *RPCMessage) {
	switch msg.Type {
	case "PING":
		client.SendPingMessage(msg)
	case "FIND_NODE":
		client.SendFindContactMessage(msg)
	case "STORE":
		client.SendStoreMessage(msg)
	case "FIND_VALUE":
		client.SendFindValueMessage(msg)
	default:
		fmt.Println("Unknown RPC message type:", msg.Type)
	}
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
*/
