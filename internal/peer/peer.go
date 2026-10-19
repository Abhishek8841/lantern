package peer

type Connection interface {
	SendMessage(peerID string, message string) error
	SendACK(peerID string, originalMessageID string) error
	Close() error
}

type Peer struct {
	ID   string
	Conn Connection
}
