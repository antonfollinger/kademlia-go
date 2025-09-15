package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

type Client struct {
	node    NodeAPI
	conn    *net.UDPConn
	pending sync.Map
}

type ClientAPI interface {
	SendPingMessage(target Contact) (RPCMessage, error)
}

func InitClient(node NodeAPI) (*Client, error) {

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, err
	}

	c := &Client{
		node: node,
		conn: conn,
	}

	fmt.Println("Client listening on: ", conn.LocalAddr())

	go c.listen()

	return c, nil
}

func (client *Client) listen() {
	buf := make([]byte, 4096)
	for {
		n, addr, err := client.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		fmt.Printf("Client found RPC from %v, bytes read: %d\n\n", addr, n)
		var resp RPCMessage
		if err := json.Unmarshal(buf[:n], &resp); err != nil {
			continue
		}

		if ch, ok := client.pending.Load(resp.PacketID); ok {
			ch.(chan RPCMessage) <- resp
			client.pending.Delete(resp.PacketID)
		}
	}
}

func (client *Client) SendMessage(target Contact, msg *RPCMessage) (chan RPCMessage, error) {

	// Build UDP address from Contact
	addr, err := net.ResolveUDPAddr("udp", target.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP addr: %w", err)
	}

	// Create response channel for this request
	respChan := make(chan RPCMessage, 1)
	client.pending.Store(msg.PacketID, respChan)

	// Marshal RPCMessage into JSON
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal RPCMessage: %w", err)
	}

	// Send JSON bytes
	_, err = client.conn.WriteToUDP(data, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to send UDP message: %w", err)
	}

	return respChan, nil
}

func (client *Client) SendPingMessage(target Contact) (RPCMessage, error) {
	// Build payload with my own contact
	payload := Payload{
		SourceContact: client.node.GetSelfContact(),
		TargetContact: target,
	}

	request := NewRPCMessage("PING", payload, true)
	fmt.Println("Sending ping...")
	respChan, err := client.SendMessage(target, request)
	if err != nil {
		return RPCMessage{}, err
	}

	select {
	case resp := <-respChan:
		fmt.Println("Response received")
		return resp, nil
	case <-time.After(2 * time.Second):
		return RPCMessage{}, fmt.Errorf("PING Timeout")
	}
}

/*
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

// For CLI implementation
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
