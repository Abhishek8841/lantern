package tcp

import (
	"fmt"
	"net"

	"github.com/abhishek8841/lantern/internal/protocol"
)

type PeerConn struct {
	conn net.Conn
}

func DialPeer(address string) (*PeerConn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("net.Dial() error: %w", err)
	}
	return &PeerConn{
		conn: conn,
	}, nil
}

func (p *PeerConn) SendMessage(message string) error {
	err := protocol.WriteFrame(p.conn, []byte(message))
	if err != nil {
		return fmt.Errorf("protocol.WriteFrame() error: %w", err)
	}
	return nil
}

func (p *PeerConn) Close() error {
	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("close peer connection: %w", err)
	}
	return nil
}

// func SendMessage(address string, message string) error {
// 	conn, err := net.Dial("tcp", address)
// 	if err != nil {
// 		return fmt.Errorf("connect to peer %s: %w", address, err)
// 	}
// 	defer conn.Close()

// 	// _, err = conn.Write([]byte(message))
// 	// if err != nil {
// 	// 	return fmt.Errorf("send message: %w", err)
// 	// }
// 	if err := protocol.WriteFrame(conn, []byte(message)); err != nil {
// 		return fmt.Errorf("send framed message: %w", err)
// 	}

// 	fmt.Println("Message sent to", address)
// 	return nil
// }
