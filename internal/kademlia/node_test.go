package kademlia

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockClient for testing
type MockClient struct{}

func (mc *MockClient) SendPingMessage(target Contact) (RPCMessage, error) {
	return RPCMessage{Type: "PONG"}, nil
}
func (mc *MockClient) SendFindNodeMessage(target *KademliaID, contact Contact) ([]Contact, error) {
	return []Contact{}, nil
}
func (mc *MockClient) SendStoreMessage(data []byte) (RPCMessage, error) {
	return RPCMessage{}, nil
}
func (mc *MockClient) SendFindValueMessage(hash string) (RPCMessage, error) {
	return RPCMessage{}, nil
}

func Test_InitNode_Bootstrap(t *testing.T) {
	node, err := InitNode(true, "localhost:8000", "")
	assert.NoError(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, "localhost:8000", node.RoutingTable.me.Address)
	assert.True(t, node.Id.Equals(NewKademliaID("0000000000000000000000000000000000000000")))
	// Should not have any contacts except self
	for i := 0; i < IDLength*8; i++ {
		assert.Equal(t, 0, node.RoutingTable.buckets[i].Len())
	}
}

func Test_InitNode_Peer(t *testing.T) {
	node, err := InitNode(false, "localhost:8001", "localhost:8000")
	assert.NoError(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, "localhost:8001", node.RoutingTable.me.Address)
	assert.False(t, node.Id.Equals(NewKademliaID("0000000000000000000000000000000000000000")))
	// Should have bootstrap contact
	found := false
	bootstrapID := NewKademliaID("0000000000000000000000000000000000000000")
	for i := 0; i < IDLength*8; i++ {
		bucket := node.RoutingTable.buckets[i]
		for e := bucket.list.Front(); e != nil; e = e.Next() {
			c := e.Value.(Contact)
			if c.ID.Equals(bootstrapID) {
				found = true
			}
		}
	}
	assert.True(t, found)
}

func Test_Node_SetClient(t *testing.T) {
	node, _ := InitNode(true, "localhost:8000", "")
	mockClient := &MockClient{}
	node.SetClient(mockClient)
	assert.Equal(t, mockClient, node.Client)
}

func Test_Node_GetSelfContact(t *testing.T) {
	node, _ := InitNode(true, "localhost:8000", "")
	self := node.GetSelfContact()
	assert.Equal(t, "localhost:8000", self.Address)
	assert.True(t, self.ID.Equals(node.Id))
}

func Test_Node_AddContact_Self(t *testing.T) {
	node, _ := InitNode(true, "localhost:8000", "")
	self := node.GetSelfContact()
	// Should not add self contact
	node.AddContact(self)
	bucketIndex := node.RoutingTable.getBucketIndex(self.ID)
	bucket := node.RoutingTable.buckets[bucketIndex]
	found := false
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		c := e.Value.(Contact)
		if c.ID.Equals(self.ID) && c.Address == self.Address {
			found = true
		}
	}
	assert.False(t, found)
}

func Test_Node_AddContact_FullBucket_Respond(t *testing.T) {
	node, _ := InitNode(true, "localhost:8000", "")
	node.SetClient(&MockClient{})

	contacts := make([]Contact, bucketSize)
	// Fill the 154:th bucket
	for i := 0; i < bucketSize; i++ {
		idStr := fmt.Sprintf("%038d%02x", 0, 40+i) // 38 zeros + 2 hex digits
		id := NewKademliaID(idStr)
		contact := Contact{ID: id, Address: fmt.Sprintf("1.2.3.4:%d", 8001+i)}
		contacts[i] = contact
		fmt.Println(contact)
		node.AddContact(contact)
	}
	bucket := node.RoutingTable.buckets[154]
	assert.Equal(t, bucketSize, bucket.Len())

	// Add a new contact to the same bucket, should trigger eviction
	idStr := fmt.Sprintf("%038d%02x", 0, 60) // 38 zeros + 2 hex digits
	newID := NewKademliaID(idStr)
	newContact := Contact{ID: newID, Address: "0.0.0.0:9999"}
	node.AddContact(newContact)
	// Bucket should still be full
	assert.Equal(t, bucketSize, bucket.Len())
	// New contact should be present
	found := false
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		c := e.Value.(Contact)
		if c.ID.Equals(newContact.ID) && c.Address == newContact.Address {
			found = true
		}
	}
	assert.False(t, found)
}

// MockClientNoRespond simulates ping failures
type MockClientNoRespond struct{}

func (mc *MockClientNoRespond) SendPingMessage(target Contact) (RPCMessage, error) {
	return RPCMessage{}, fmt.Errorf("no response")
}
func (mc *MockClientNoRespond) SendFindNodeMessage(target *KademliaID, contact Contact) ([]Contact, error) {
	return []Contact{}, nil
}
func (mc *MockClientNoRespond) SendStoreMessage(data []byte) (RPCMessage, error) {
	return RPCMessage{}, nil
}
func (mc *MockClientNoRespond) SendFindValueMessage(hash string) (RPCMessage, error) {
	return RPCMessage{}, nil
}

func Test_Node_AddContact_FullBucket_NoRespond(t *testing.T) {
	node, _ := InitNode(true, "localhost:8000", "")
	node.SetClient(&MockClientNoRespond{})

	contacts := make([]Contact, bucketSize)
	// Fill the 154:th bucket
	for i := 0; i < bucketSize; i++ {
		idStr := fmt.Sprintf("%038d%02x", 0, 40+i) // 38 zeros + 2 hex digits
		id := NewKademliaID(idStr)
		contact := Contact{ID: id, Address: fmt.Sprintf("1.2.3.4:%d", 8001+i)}
		contacts[i] = contact
		node.AddContact(contact)
	}
	bucket := node.RoutingTable.buckets[154]
	assert.Equal(t, bucketSize, bucket.Len())

	// Add a new contact to the same bucket, should trigger eviction
	idStr := fmt.Sprintf("%038d%02x", 0, 60) // 38 zeros + 2 hex digits
	newID := NewKademliaID(idStr)
	newContact := Contact{ID: newID, Address: "0.0.0.0:9999"}
	node.AddContact(newContact)
	// Bucket should still be full
	assert.Equal(t, bucketSize, bucket.Len())
	// New contact should be present
	found := false
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		c := e.Value.(Contact)
		if c.ID.Equals(newContact.ID) && c.Address == newContact.Address {
			found = true
		}
	}
	assert.True(t, found)
}

func Test_Node_AddContact_Duplicate(t *testing.T) {
	node, _ := InitNode(true, "localhost:8000", "")
	id := NewKademliaID("0000000000000000000000000000000000000001")
	contact := Contact{ID: id, Address: "localhost:8001"}
	node.AddContact(contact)
	node.AddContact(contact) // Add duplicate
	bucketIndex := node.RoutingTable.getBucketIndex(contact.ID)
	bucket := node.RoutingTable.buckets[bucketIndex]
	count := 0
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		c := e.Value.(Contact)
		if c.ID.Equals(contact.ID) && c.Address == contact.Address {
			count++
		}
	}
	assert.Equal(t, 1, count)
}

func Test_Node_AddContact_Normal(t *testing.T) {
	node, _ := InitNode(true, "localhost:8000", "")
	id := NewKademliaID("0000000000000000000000000000000000000001")
	contact := Contact{ID: id, Address: "localhost:8001"}
	node.AddContact(contact)
	found := false
	bucketIndex := node.RoutingTable.getBucketIndex(contact.ID)
	bucket := node.RoutingTable.buckets[bucketIndex]
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		c := e.Value.(Contact)
		if c.ID.Equals(contact.ID) && c.Address == contact.Address {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func Test_Node_LookupClosestContacts(t *testing.T) {
	node, _ := InitNode(true, "localhost:8000", "")
	for i := 0; i < 5; i++ {
		id := NewRandomKademliaID()
		contact := Contact{ID: id, Address: "localhost:8001"}
		node.AddContact(contact)
	}
	target := Contact{ID: NewRandomKademliaID(), Address: "localhost:8002"}
	closest := node.LookupClosestContacts(target)
	assert.LessOrEqual(t, len(closest), alpha)
	for _, c := range closest {
		assert.NotNil(t, c.ID)
	}
}

func Test_Node_IterativeFindNode(t *testing.T) {
	node, _ := InitNode(true, "localhost:8000", "")
	node.SetClient(&MockClient{})

	// Add multiple contacts to the node
	numContacts := 10
	for i := 0; i < numContacts; i++ {
		id := NewRandomKademliaID()
		contact := Contact{ID: id, Address: fmt.Sprintf("localhost:%d", 8001+i)}
		node.AddContact(contact)
	}

	target := NewRandomKademliaID()
	contacts, err := node.IterativeFindNode(target)
	assert.NoError(t, err)
	assert.NotNil(t, contacts)
	// Should not contain self contact
	for _, c := range contacts {
		assert.False(t, c.ID.Equals(node.Id))
	}
}

func Test_Node_LookupData_Store(t *testing.T) {
	node, _ := InitNode(true, "localhost:8000", "")
	key := "testkey"
	data := []byte("testdata")
	node.Store(key, data)
	result := node.LookupData(key)
	assert.Equal(t, data, result)
	// Lookup for non-existent key
	result2 := node.LookupData("notfound")
	assert.Nil(t, result2)
}

func Test_Node_PrintStore(t *testing.T) {
	node, _ := InitNode(true, "localhost:8000", "")
	node.Store("key", []byte("value"))
	node.PrintStore() // Just ensure no panic
}

func Test_Node_PrintRoutingTable(t *testing.T) {
	node, _ := InitNode(true, "localhost:8000", "")
	node.PrintRoutingTable() // Just ensure no panic
}
