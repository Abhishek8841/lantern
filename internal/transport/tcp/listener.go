package tcp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/abhishek8841/lantern/internal/peer"
	"github.com/abhishek8841/lantern/internal/protocol"
)

// fix the unhandled goroutine leak here is .Accept() errors ???

func StartListener(ctx context.Context, address string, manager *peer.Manager, ID string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		// this wraps the err inside another error ??? -> error.is() gives true in the wrapping case but == doesn't
		return fmt.Errorf("start TCP listener: %w", err)
	}
	defer listener.Close()

	var handlers sync.WaitGroup

	shutdown := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = listener.Close()
		case <-shutdown:
			return
		}
	}()
	defer close(shutdown)

	fmt.Println("LANtern listening on", listener.Addr())

	var returnError error = nil

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			returnError = fmt.Errorf("accept connection: %w", err)
			break
		}

		if ctx.Err() != nil {
			_ = conn.Close()
			break
		}

		fmt.Println("Peer connected:", conn.RemoteAddr())

		handlers.Add(1)

		go func(conn net.Conn) {
			defer handlers.Done()
			serverConnection(conn, ID, ctx, manager)
		}(conn)
	}
	_ = manager.CloseAll()
	handlers.Wait()
	return returnError
}

func serverConnection(conn net.Conn, ID string, ctx context.Context, manager *peer.Manager) {
	defer conn.Close()
	pConn := NewPeerConn(conn)
	remoteID, err := performHandshake(pConn, ID)
	if err != nil {
		if ctx.Err() == nil {
			fmt.Println("handshake failed:", err)
		}
		return
	}

	p := &peer.Peer{
		ID:   remoteID,
		Conn: conn,
	}

	successfullAdd := manager.Add(p)
	if !successfullAdd {
		_ = conn.Close()
		return
	}

	defer manager.Remove(p)
	if err := handleConnection(pConn, ID, remoteID); err != nil {
		if ctx.Err() == nil {
			fmt.Println("Connection ended:", err)
		}
		return
	}

}

// func handleConnection(conn net.Conn) error {
// 	defer conn.Close()

// 	buffer := make([]byte, 1024)

// 	for {
// 		n, err := conn.Read(buffer)

// 		if n > 0 {
// 			fmt.Printf("Received: %s from %s\n", string(buffer[:n]), conn.RemoteAddr())
// 		}

// 		if err != nil {
// 			if err == io.EOF {
// 				fmt.Println("Peer disconnected")
// 				return nil
// 			}
// 			return fmt.Errorf("read from peer: %w", err)
// 		}
// 	}
// }

func performHandshake(p *PeerConn, localID string) (string, error) {
	remoteID, err := p.Handshake(localID)
	if err != nil {
		return "", fmt.Errorf("handshake with %s: %w", p.conn.RemoteAddr(), err)
	}
	return remoteID, nil
}

func handleConnection(p *PeerConn, localID, remoteID string) error {
	// handshake, err := protocol.ReadFrame(conn)
	// if err != nil {
	// 	if errors.Is(err, io.EOF) {
	// 		return nil
	// 	}
	// 	return fmt.Errorf("read HELLO message from %s: %w", conn.RemoteAddr(), err)
	// }
	// envelope, err := protocol.DecodeMessage(handshake)
	// if err != nil {
	// 	return fmt.Errorf("decode HELLO message from %s: %w", conn.RemoteAddr(), err)
	// }
	// if envelope.Type != protocol.MessageTypeHello {
	// 	return fmt.Errorf("expected HELLO as first message")
	// }
	// if envelope.SenderID == "" {
	// 	return fmt.Errorf("HELLO has empty sender ID")
	// }
	// peerID := envelope.SenderID
	// fmt.Println("Peer identified:", peerID)

	for {
		envelope, err := p.ReadEnvelope()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Peer disconnected:", p.conn.RemoteAddr())
				return nil
			}
			return fmt.Errorf("read message from %s: %w", p.conn.RemoteAddr(), err)
		}

		switch envelope.Type {
		case protocol.MessageTypeMessage:
			if envelope.SenderID != remoteID {
				return fmt.Errorf(
					"sender ID mismatch: expected %q, got %q",
					remoteID,
					envelope.SenderID,
				)
			}

			if envelope.ID == "" {
				return fmt.Errorf("received MESSAGE with empty ID")
			}

			if err := p.SendACK(localID, envelope.ID); err != nil {
				return fmt.Errorf("send ACK for message %s: %w", envelope.ID, err)
			}

			fmt.Printf("[%s] %s\n", remoteID, envelope.Payload)
		default:
			return fmt.Errorf("unsupported message type: %q", envelope.Type)
		}
	}
}
