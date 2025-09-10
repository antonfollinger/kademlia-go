package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_contact_NewContact(t *testing.T) {
	id := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	address := "127.0.0.1"
	contact := NewContact(id, address)

	assert.Equal(t, id, contact.ID)
	assert.Equal(t, address, contact.Address)
	assert.Nil(t, contact.distance)
}

func Test_contact_CalcDistance(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	id2 := NewKademliaID("0000000000000000000000000000000000000001")
	contact := NewContact(id1, "127.0.0.1")
	contact.CalcDistance(id2)

	assert.NotNil(t, contact.distance)
	expected := id1.CalcDistance(id2)
	assert.Equal(t, expected, contact.distance)
}

func Test_contact_Less(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	id2 := NewKademliaID("0000000000000000000000000000000000000001")
	id3 := NewKademliaID("0000000000000000000000000000000000000002")

	contact1 := NewContact(id1, "127.0.0.1")
	contact2 := NewContact(id2, "127.0.0.2")
	contact3 := NewContact(id3, "127.0.0.3")

	target := NewKademliaID("0000000000000000000000000000000000000000")
	contact1.CalcDistance(target)
	contact2.CalcDistance(target)
	contact3.CalcDistance(target)

	assert.False(t, contact1.Less(&contact2))
	assert.True(t, contact2.Less(&contact3))
}

func Test_contact_String(t *testing.T) {
	id := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	contact := NewContact(id, "127.0.0.1")
	// KademliaID.string() produces lowercase hexadecimal string
	expected := `contact("ffffffffffffffffffffffffffffffffffffffff", "127.0.0.1")`
	assert.Equal(t, expected, contact.String())
}

func Test_contactcandidates_AppendAndGetContacts(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	id2 := NewKademliaID("0000000000000000000000000000000000000001")
	contact1 := NewContact(id1, "127.0.0.1")
	contact2 := NewContact(id2, "127.0.0.2")

	candidates := &ContactCandidates{}
	candidates.Append([]Contact{contact1, contact2})

	assert.Equal(t, 2, candidates.Len())
	contacts := candidates.GetContacts(1)
	assert.Equal(t, 1, len(contacts))
	assert.Equal(t, contact1, contacts[0])
}

func Test_contactcandidates_Swap(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	id2 := NewKademliaID("0000000000000000000000000000000000000001")
	contact1 := NewContact(id1, "127.0.0.1")
	contact2 := NewContact(id2, "127.0.0.2")

	candidates := &ContactCandidates{}
	candidates.Append([]Contact{contact1, contact2})
	candidates.Swap(0, 1)

	assert.Equal(t, contact2, candidates.contacts[0])
	assert.Equal(t, contact1, candidates.contacts[1])
}

func Test_contactcandidates_Less(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	id2 := NewKademliaID("0000000000000000000000000000000000000001")
	contact1 := NewContact(id1, "127.0.0.1")
	contact2 := NewContact(id2, "127.0.0.2")

	target := NewKademliaID("0000000000000000000000000000000000000000")
	contact1.CalcDistance(target)
	contact2.CalcDistance(target)

	candidates := &ContactCandidates{}
	candidates.Append([]Contact{contact1, contact2})

	assert.False(t, candidates.Less(0, 1))
	assert.True(t, candidates.Less(1, 0))
}

func Test_contactcandidates_Sort(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	id2 := NewKademliaID("0000000000000000000000000000000000000001")
	id3 := NewKademliaID("0000000000000000000000000000000000000002")
	contact1 := NewContact(id1, "127.0.0.1")
	contact2 := NewContact(id2, "127.0.0.2")
	contact3 := NewContact(id3, "127.0.0.3")

	target := NewKademliaID("0000000000000000000000000000000000000000")
	contact1.CalcDistance(target)
	contact2.CalcDistance(target)
	contact3.CalcDistance(target)

	candidates := &ContactCandidates{}
	candidates.Append([]Contact{contact1, contact3, contact2})
	candidates.Sort()

	assert.Equal(t, contact2, candidates.contacts[0])
	assert.Equal(t, contact3, candidates.contacts[1])
	assert.Equal(t, contact1, candidates.contacts[2])
}
