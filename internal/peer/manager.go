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
	// concurrency -> prevent 2 connection check
	if _, exists := m.peers[p.ID]; exists {
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

func (m *Manager) Remove(p *Peer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	current, ok := m.peers[p.ID]
	if ok && current == p {
		delete(m.peers, p.ID)
	}
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

func (m *Manager) SendMessage(localID, peerID, message string) error {
	m.mu.RLock()
	p, ok := m.peers[peerID]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("peer %q is not connected", peerID)
	}
	return p.Conn.SendMessage(localID, message)
}

func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ids := make([]string, 0, len(m.peers))
	for id := range m.peers {
		ids = append(ids, id)
	}
	return ids
}
