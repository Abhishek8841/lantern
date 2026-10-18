package peer

import (
	"fmt"
	"sync"
)

type Manager struct {
	mu       sync.RWMutex
	peers    map[string]*Peer
	shutting bool
}

func NewManager() *Manager {
	return &Manager{
		peers: make(map[string]*Peer),
	}
}

func (m *Manager) Add(p *Peer) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shutting {
		return false
	}

	m.peers[p.ID] = p
	return true
}

func (m *Manager) Get(id string) (*Peer, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.peers[id]
	return p, ok
}

func (m *Manager) Remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.peers, id)
}

func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.peers)
}

func (m *Manager) CloseAll() error {
	m.mu.Lock()

	m.shutting = true

	peers := make([]*Peer, 0, len(m.peers))
	for _, p := range m.peers {
		peers = append(peers, p)
	}

	m.mu.Unlock()

	var firstError error

	for _, p := range peers {
		if err := p.Conn.Close(); err != nil && firstError == nil {
			firstError = fmt.Errorf("close peer %s: %w", p.ID, err)
		}
	}
	return firstError
}

// select {} quick hack to keep a goroutine alive to test temporarily
