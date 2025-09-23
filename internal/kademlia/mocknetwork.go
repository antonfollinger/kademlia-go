package kademlia

import (
	"fmt"
	"sync"
)

type mockPacket struct {
	src  string
	data []byte
}

type MockRegistry struct {
	channels map[string]chan mockPacket
	closed   map[string]bool
	mu       sync.RWMutex
}

type MockNetwork struct {
	addr     string
	registry *MockRegistry
}

func NewMockNetwork(addr string, registry *MockRegistry) *MockNetwork {
	registry.Register(addr)
	return &MockNetwork{
		addr:     addr,
		registry: registry,
	}
}

func NewMockRegistry() *MockRegistry {
	return &MockRegistry{channels: make(map[string]chan mockPacket), closed: make(map[string]bool)}
}

func (r *MockRegistry) Register(addr string) chan mockPacket {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch := make(chan mockPacket, 50)
	r.channels[addr] = ch
	r.closed[addr] = false
	return ch
}

func (r *MockRegistry) Get(addr string) (chan mockPacket, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ch, ok := r.channels[addr]
	if !ok || r.closed[addr] {
		return nil, false
	}
	return ch, ok
}

func (m *MockNetwork) GetConn() string {
	return m.addr
}

func (m *MockNetwork) SendMessage(addr string, data []byte) error {
	ch, ok := m.registry.Get(addr)
	if !ok {
		return fmt.Errorf("address not found or channel closed: %s", addr)
	}
	select {
	case ch <- mockPacket{src: m.addr, data: data}:
		return nil
	default:
		return fmt.Errorf("channel full or closed for address: %s", addr)
	}
}

func (m *MockNetwork) ReceiveMessage() (string, []byte, error) {
	ch, ok := m.registry.Get(m.addr)
	if !ok {
		return "", nil, fmt.Errorf("no channel for address: %s", m.addr)
	}
	pkt, ok := <-ch
	if !ok {
		return "", nil, fmt.Errorf("channel closed for address: %s", m.addr)
	}
	return pkt.src, pkt.data, nil
}

func (m *MockNetwork) Close() error {
	m.registry.mu.Lock()
	ch, ok := m.registry.channels[m.addr]
	if ok && !m.registry.closed[m.addr] {
		close(ch)
		m.registry.closed[m.addr] = true
	}
	m.registry.mu.Unlock()
	return nil
}
