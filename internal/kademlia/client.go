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

// When part of a network, it must be possible for any node to upload an object
// that will end up at the designated storage nodes. In Kademlia terminology,
// the designated nodes are the K nodes nearest to the hash of the data object in question.
// Data objects are always UTF-8 strings
func (client *Client) SendStoreMessage(data []byte) (RPCMessage, error) {
	// Use a hashing method to generate a KademliaID key from the data
	key := NewKademliaID("abc") // change to data here

	// Find closest nodes to the generated key
	closest, err := client.node.IterativeFindNode(key)
	if err != nil {
		return RPCMessage{}, err
	}

	k := alpha // minimum number of nodes to store
	storedCount := 0
	var lastResp RPCMessage

	for _, contact := range closest {
		request := NewRPCMessage("STORE", Payload{Key: key.String(), Data: data}, true)
		respChan, err := client.SendMessage(contact, request)
		if err != nil {
			continue
		}

		select {
		case resp := <-respChan:
			fmt.Println("STORE response received")

			// Assume any response means successful store
			storedCount++
			lastResp = resp
			if storedCount >= k {
				return lastResp, nil
			}
		case <-time.After(2 * time.Second):
			fmt.Println("STORE Timeout for contact", contact.String())
			// try next contact
		}
	}
	// If fewer than k nodes stored the data
	return RPCMessage{}, fmt.Errorf("data could not be stored on at least %d nodes", k)
}

// When part of a network with uploaded objects, it must be possible to find and
// download any object, as long as it is stored by at least one designated node.
func (client *Client) SendFindValueMessage(hash string) (RPCMessage, error) {
	hashID := NewKademliaID(hash)
	closest, err := client.node.IterativeFindNode(hashID)
	if err != nil {
		return RPCMessage{}, err
	}

	for _, contact := range closest {
		request := NewRPCMessage("FIND_VALUE", Payload{Key: hash}, true)
		respChan, err := client.SendMessage(contact, request)
		if err != nil {
			continue
		}

		select {
		case resp := <-respChan:
			fmt.Println("FIND_VALUE response received")
			if resp.Payload.Data != nil {
				// Found the data, return immediately
				return resp, nil
			}
			// else, try next contact
		case <-time.After(2 * time.Second):
			fmt.Println("FIND_VALUE Timeout for contact", contact.String())
			// try next contact
		}
	}
	// If none of the contacts had the data
	return RPCMessage{}, fmt.Errorf("FIND_VALUE not found on any contacted node")
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
