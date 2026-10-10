package tcp

import (
	"fmt"
	"io"
	"net"
)

func StartListener(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("start TCP listener: %w", err)
	}
	defer listener.Close()

	fmt.Println("LANtern listening on", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("accept connection: %w", err)
		}

		fmt.Println("Peer connected:", conn.RemoteAddr())

		go func(connection net.Conn) {
			if err := handleConnection(conn); err != nil {
				fmt.Println("Connection ended:", err)
			}
		}(conn)
	}
}

func handleConnection(conn net.Conn) error {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)

		if n > 0 {
			fmt.Printf("Received: %s from %s\n", string(buffer[:n]), conn.RemoteAddr())
		}

		if err != nil {
			if err == io.EOF {
				fmt.Println("Peer disconnected")
				return nil
			}
			return fmt.Errorf("read from peer: %w", err)
		}
	}
}
