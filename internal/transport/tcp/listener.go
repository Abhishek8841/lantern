package tcp

import (
	"fmt"
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
		conn.Close()
	}
}
