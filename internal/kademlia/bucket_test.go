package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	bucket := newBucket()

	assert.True(t, bucket.Len() == 0)

	contact := NewContact(NewRandomKademliaID(), "0.0.0.0:1234")

	bucket.AddContact(contact)

	assert.True(t, bucket.Len() == 1)
}

func TestAddContact(t *testing.T) {
	bucket := newBucket()

	// Add c1 to bucket and confirm it's at the front of the list
	contact1 := NewContact(NewRandomKademliaID(), "0.0.0.0:1234")
	bucket.AddContact(contact1)
	assert.True(t, bucket.list.Front().Value == contact1)

	// Now add c2 to bucket and confirm it is now at the front of the list
	contact2 := NewContact(NewRandomKademliaID(), "0.0.0.0:5678")
	bucket.AddContact(contact2)
	assert.True(t, bucket.list.Front().Value == contact2)

	// Add c1 to bucket again and confirm its back at the front of the list
	bucket.AddContact(contact1)
	assert.True(t, bucket.list.Front().Value == contact1)
}
