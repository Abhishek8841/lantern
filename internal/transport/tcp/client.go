package tcp

import (
	"fmt"
	"net"
	"time"

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

func (p *PeerConn) Close() error {
	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("close peer connection: %w", err)
	}
	return nil
}

func (p *PeerConn) sendEnvelope(message protocol.Envelope) error {
	data, err := protocol.EncodeMessage(message)
	if err != nil {
		return fmt.Errorf("encode message: %w", err)
	}
	if err := protocol.WriteFrame(p.conn, data); err != nil {
		return fmt.Errorf("protocol.WriteFrame() error: %w", err)
	}
	return nil
}

func (p *PeerConn) SendHello(peerID string) error {
	return p.sendEnvelope(protocol.Envelope{
		Type:      protocol.MessageTypeHello,
		ID:        protocol.NewMessageID(),
		SenderID:  peerID,
		Timestamp: time.Now().UTC(),
	})
}

func (p *PeerConn) SendMessage(peerID string, message string) error {
	return p.sendEnvelope(protocol.Envelope{
		Type:      protocol.MessageTypeMessage,
		ID:        protocol.NewMessageID(),
		SenderID:  peerID,
		Timestamp: time.Now().UTC(),
		Payload:   message,
	})
}

func (p *PeerConn) SendACK(peerID string, originalMessageID string) error {
	return p.sendEnvelope(protocol.Envelope{
		Type:           protocol.MessageTypeAck,
		ID:             protocol.NewMessageID(),
		SenderID:       peerID,
		Timestamp:      time.Now().UTC(),
		AcknowledgesID: originalMessageID,
	})
}

func (p *PeerConn) ReadACKFrame() (protocol.Envelope, error) {
	data, err := protocol.ReadFrame(p.conn)
	if err != nil {
		return protocol.Envelope{}, fmt.Errorf("read ACK frame error: %w", err)
	}
	message, err := protocol.DecodeMessage(data)
	if err != nil {
		return protocol.Envelope{}, fmt.Errorf("decode ACK envelope error: %w", err)
	}
	if message.Type != protocol.MessageTypeAck {
		return protocol.Envelope{}, fmt.Errorf("expected ACK, received %q", message.Type)
	}

	return message, nil
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
