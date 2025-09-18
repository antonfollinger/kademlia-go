package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_routingtable_GetMe(t *testing.T) {
	id := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	contact := Contact{ID: id, Address: "localhost:8000"}
	rt := NewRoutingTable(contact)

	me := GetMe(rt)

	assert.Equal(t, me, contact)
}

func Test_routingtable_NewRoutingTable(t *testing.T) {
	id := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	contact := Contact{ID: id, Address: "localhost:8000"}
	rt := NewRoutingTable(contact)

	assert.NotNil(t, rt)
	assert.Equal(t, contact, rt.me)
	assert.Len(t, rt.buckets, IDLength*8)
	for i := 0; i < IDLength*8; i++ {
		assert.NotNil(t, rt.buckets[i])
	}
}

func Test_routingtable_AddContact(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	id2 := NewKademliaID("0000000000000000000000000000000000000001")
	contact1 := Contact{ID: id1, Address: "localhost:8000"}
	contact2 := Contact{ID: id2, Address: "localhost:8001"}
	rt := NewRoutingTable(contact1)

	rt.AddContact(contact2)
	bucketIndex := rt.getBucketIndex(contact2.ID)
	bucket := rt.buckets[bucketIndex]

	found := false
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		c := e.Value.(Contact)
		if c.ID.Equals(contact2.ID) && c.Address == contact2.Address {
			found = true
			break
		}
	}

	assert.True(t, found)
}

func Test_routingtable_FindClosestContacts(t *testing.T) {
	id := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	rt := NewRoutingTable(Contact{ID: id, Address: "localhost:8000"})

	// Add several contacts
	for i := 0; i < 10; i++ {
		contactID := NewRandomKademliaID()
		contact := Contact{ID: contactID, Address: "localhost:8000"}
		rt.AddContact(contact)
	}

	target := NewKademliaID("0000000000000000000000000000000000000001")
	closest := rt.FindClosestContacts(target, 5)
	assert.LessOrEqual(t, len(closest), 5)
	for _, c := range closest {
		assert.NotNil(t, c.ID)
	}
}

func Test_routingtable_getBucketIndex(t *testing.T) {
	me := Contact{ID: NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")}
	rt := NewRoutingTable(me)

	id := NewKademliaID("0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	index := rt.getBucketIndex(id)
	assert.Equal(t, index, 0)
}

func Test_routingtable_getBucketIndex_self(t *testing.T) {
	me := Contact{ID: NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")}
	rt := NewRoutingTable(me)

	index := rt.getBucketIndex(me.ID)
	assert.Equal(t, index, 159)
}

func Test_routingtable_FindClosestContacts_EmptyTable(t *testing.T) {
	id := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	rt := NewRoutingTable(Contact{ID: id, Address: "localhost:8000"})
	target := NewKademliaID("0000000000000000000000000000000000000001")
	closest := rt.FindClosestContacts(target, 3)
	assert.Equal(t, 0, len(closest))
}
