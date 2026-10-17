package peer

import (
	"fmt"
	"sync"
)

type Manager struct {
	mu    sync.RWMutex
	peers map[string]*Peer
}

func NewManager() *Manager {
	return &Manager{
		peers: make(map[string]*Peer),
	}
}

func (m *Manager) Add(p *Peer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.peers[p.ID] = p
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
	defer m.mu.Unlock()

	var firstError error

	for id, p := range m.peers {
		if err := p.Conn.Close(); err != nil && firstError == nil {
			firstError = fmt.Errorf("close peer %s: %w", id, err)
		}
	}
	return firstError
}
