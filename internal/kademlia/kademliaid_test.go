package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_KademliaID_NewKademliaID(t *testing.T) {

	// 20 bytes
	id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewKademliaID("1234567891234567891234567891234567891234")

	assert.Equal(t, id1.String(), "ffffffff00000000000000000000000000000000")
	assert.Equal(t, id2.String(), "1234567891234567891234567891234567891234")
}

func Test_KademliaID_NewRandomKademliaID(t *testing.T) {
	id1 := NewRandomKademliaID()
	id2 := NewRandomKademliaID()

	assert.True(t, id1.Less(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")))
	assert.False(t, id1.Equals(id2))
}

func Test_KademliaID_CalcDistance(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewKademliaID("5555555500000000000000000000000000000000")

	// F XOR 5 = A
	d := NewKademliaID("AAAAAAAA000000000000000000000000000000000")
	assert.Equal(t, id1.CalcDistance(id2), d)
}

func Test_KademliaID_Equals(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewKademliaID("1234567891234567891234567891234567891234")
	id3 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")

	assert.True(t, id1.Equals(id3))
	assert.False(t, id1.Equals(id2))
}
func Test_KademliaID_Less(t *testing.T) {
	// id1 < id2
	id1 := NewKademliaID("0000000000000000000000000000000000000000")
	id2 := NewKademliaID("0000000000000000000000000000000000000001")
	assert.True(t, id1.Less(id2))
	assert.False(t, id2.Less(id1))

	// id1 == id2
	id3 := NewKademliaID("ABCDEFABCDEFABCDEFABCDEFABCDEFABCDEFABCD")
	id4 := NewKademliaID("ABCDEFABCDEFABCDEFABCDEFABCDEFABCDEFABCD")
	assert.False(t, id3.Less(id4))
	assert.False(t, id4.Less(id3))

	// id1 > id2
	id5 := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	id6 := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFE")
	assert.False(t, id5.Less(id6))
	assert.True(t, id6.Less(id5))
}
