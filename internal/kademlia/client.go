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
	SendFindNodeMessage(target *KademliaID, contact Contact) ([]Contact, error)
}

func InitClient(node NodeAPI) (*Client, error) {

	// Create connection with ephemeral port
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

	// Add contact information
	msg.Payload.SourceContact = client.node.GetSelfContact()
	msg.Payload.TargetContact = target

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

	request := NewRPCMessage("PING", Payload{}, true)
	respChan, err := client.SendMessage(target, request)
	if err != nil {
		return RPCMessage{}, err
	}

	// Wait for response
	select {
	case resp := <-respChan:
		// Add contact
		fmt.Println("PING response received")
		return resp, nil
	case <-time.After(2 * time.Second):
		return RPCMessage{}, fmt.Errorf("PING Timeout")
	}
}

// JOIN, PING BOOTSTRAP, FIND_NODE SELF -> UNTIL DISTANCE ISN'T GETTING SMALLER
// BASE CASE DISTANCE TO NODE AND TARGET IS 0
func (client *Client) SendFindNodeMessage(target *KademliaID, contact Contact) ([]Contact, error) {

	payload := Payload{
		Key: target.String(),
	}

	request := NewRPCMessage("FIND_NODE", payload, true)
	respChan, err := client.SendMessage(contact, request)
	if err != nil {
		return nil, err
	}

	// Wait for response
	select {
	case resp := <-respChan:
		if resp.Payload.SourceContact != (Contact{}) {
			client.node.AddContact(resp.Payload.SourceContact)
		}
		fmt.Println("FIND_NODE response received")
		for _, c := range resp.Payload.Contacts {
			client.node.AddContact(c)
		}
		return resp.Payload.Contacts, nil
	case <-time.After(2 * time.Second):
		return nil, fmt.Errorf("FIND_NODE Timeout")
	}
}

/*
func (client *Client) SendStoreMessage(target Contact) (RPCMessage, error) {

	request := NewRPCMessage("STORE", Payload{}, true)
	respChan, err := client.SendMessage(target, request)
	if err != nil {
		return RPCMessage{}, err
	}

	// Wait for response
	select {
	case resp := <-respChan:
		fmt.Println("Response received")
		fmt.Printf("RPC INFO: %+v\n\n", resp)
		return resp, nil
	case <-time.After(2 * time.Second):
		return RPCMessage{}, fmt.Errorf("STORE Timeout")
	}

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

func (client *Client) SendFindValueMessage(target Contact) (RPCMessage, error) {

	request := NewRPCMessage("FIND_VALUE", Payload{}, true)
	respChan, err := client.SendMessage(target, request)
	if err != nil {
		return RPCMessage{}, err
	}

	// Wait for response
	select {
	case resp := <-respChan:
		fmt.Println("Response received")
		fmt.Printf("RPC INFO: %+v\n\n", resp)
		return resp, nil
	case <-time.After(2 * time.Second):
		return RPCMessage{}, fmt.Errorf("FIND_VALUE Timeout")
	}

	/*
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
/*
func (client *Client) handleRPC(msg *RPCMessage) {
	switch msg.Type {
	case "PING":
		client.SendPingMessage(msg)
	case "FIND_NODE":
		client.SendFindNodeMessage(msg)
	case "STORE":
		client.SendStoreMessage(msg)
	case "FIND_VALUE":
		client.SendFindValueMessage(msg)
	default:
		fmt.Println("Unknown RPC message type:", msg.Type)
	}
}
*/
