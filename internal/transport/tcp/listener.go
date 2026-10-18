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

func StartListener(ctx context.Context, address string, manager *peer.Manager) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		// this wraps the err inside another error ??? -> error.is() gives true in the wrapping case but == doesn't
		return fmt.Errorf("start TCP listener: %w", err)
	}
	defer listener.Close()

	var handlers sync.WaitGroup

	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

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

		p := &peer.Peer{
			ID:   conn.RemoteAddr().String(),
			Conn: conn,
		}

		successfullAdd := manager.Add(p)
		if !successfullAdd {
			_ = conn.Close()
			continue
		}
		handlers.Add(1)

		fmt.Println("Peer connected:", conn.RemoteAddr())

		go func(p *peer.Peer) {
			defer handlers.Done()
			defer p.Conn.Close()
			defer manager.Remove(p.ID)
			if err := handleConnection(p.Conn); err != nil {
				if ctx.Err() != nil {
					return
				}
				fmt.Println("Connection ended:", err)
			}
		}(p)
	}
	_ = manager.CloseAll()
	handlers.Wait()
	return returnError
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

func handleConnection(conn net.Conn) error {

	handshake, err := protocol.ReadFrame(conn)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return fmt.Errorf("read HELLO message from %s: %w", conn.RemoteAddr(), err)
	}
	envelope, err := protocol.DecodeMessage(handshake)
	if err != nil {
		return fmt.Errorf("decode HELLO message from %s: %w", conn.RemoteAddr(), err)
	}
	if envelope.Type != protocol.MessageTypeHello {
		return fmt.Errorf("expected HELLO as first message")
	}
	if envelope.SenderID == "" {
		return fmt.Errorf("HELLO has empty sender ID")
	}

	peerID := envelope.SenderID
	fmt.Println("Peer identified:", peerID)

	for {
		message, err := protocol.ReadFrame(conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Peer disconnected:", conn.RemoteAddr())
				return nil
			}
			return fmt.Errorf("read message from %s: %w", conn.RemoteAddr(), err)
		}

		envelope, err := protocol.DecodeMessage(message)
		if err != nil {
			return fmt.Errorf("decode message from %s: %w", conn.RemoteAddr(), err)
		}

		switch envelope.Type {
		case protocol.MessageTypeMessage:
			if envelope.SenderID != peerID {
				return fmt.Errorf(
					"sender ID mismatch: expected %q, got %q",
					peerID,
					envelope.SenderID,
				)
			}
			fmt.Printf("[%s] %s\n", peerID, envelope.Payload)
		default:
			return fmt.Errorf("unsupported message type: %q", envelope.Type)
		}
	}
}
