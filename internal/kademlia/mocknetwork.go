package kademlia

import (
	"fmt"
	"math/rand"
	"sync"
)

type mockPacket struct {
	src  string
	data []byte
}

type MockRegistry struct {
	channels map[string]chan mockPacket
	mu       sync.RWMutex
}

type MockNetwork struct {
	addr     string
	registry *MockRegistry
	dropRate float64
}

func NewMockNetwork(addr string, registry *MockRegistry, droprate float64) *MockNetwork {
	registry.Register(addr)
	return &MockNetwork{
		addr:     addr,
		registry: registry,
		dropRate: droprate,
	}
}

func NewMockRegistry() *MockRegistry {
	return &MockRegistry{channels: make(map[string]chan mockPacket)}
}

func (r *MockRegistry) Register(addr string) chan mockPacket {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch := make(chan mockPacket, 1000)
	r.channels[addr] = ch
	return ch
}

func (r *MockRegistry) Get(addr string) (chan mockPacket, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ch, ok := r.channels[addr]
	return ch, ok
}

func (m *MockNetwork) GetConn() string {
	return m.addr
}

func (m *MockNetwork) SendMessage(addr string, data []byte) error {
	if rand.Float64() < m.dropRate {
		return nil // simulate drop
	}
	if ch, ok := m.registry.Get(addr); ok {
		ch <- mockPacket{src: m.addr, data: data}
		return nil
	}
	return fmt.Errorf("address not found: %s", addr)
}

func (m *MockNetwork) ReceiveMessage() (string, []byte, error) {
	ch, ok := m.registry.Get(m.addr)
	if !ok {
		return "", nil, fmt.Errorf("no channel for address: %s", m.addr)
	}
	pkt := <-ch
	return pkt.src, pkt.data, nil
}

func (m *MockNetwork) Close() error {
	ch, ok := m.registry.Get(m.addr)
	if ok {
		close(ch)
	}
	return nil
}
