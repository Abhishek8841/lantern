package tcp

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/abhishek8841/lantern/internal/peer"
	"github.com/abhishek8841/lantern/internal/protocol"
)

func StartListener(address string, manager *peer.Manager) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		// this wraps the err inside another error ??? -> error.is() gives true in the wrapping case but == doesn't
		return fmt.Errorf("start TCP listener: %w", err)
	}
	defer listener.Close()

	fmt.Println("LANtern listening on", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("accept connection: %w", err)
		}

		p := &peer.Peer{
			ID:   conn.RemoteAddr().String(),
			Conn: conn,
		}

		manager.Add(p)
		fmt.Println("Peer connected:", conn.RemoteAddr())

		go func(p *peer.Peer) {
			defer p.Conn.Close()
			defer manager.Remove(p.ID)
			if err := handleConnection(p.Conn); err != nil {
				fmt.Println("Connection ended:", err)
			}
		}(p)
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

func handleConnection(conn net.Conn) error {
	for {
		message, err := protocol.ReadFrame(conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Peer disconnected:", conn.RemoteAddr())
				return nil
			}
			return fmt.Errorf("read message from %s: %w", conn.RemoteAddr(), err)
		}

		fmt.Printf("Received from %s: %s\n", conn.RemoteAddr(), string(message))
	}
}
