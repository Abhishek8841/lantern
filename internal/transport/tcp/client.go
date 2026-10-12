package tcp

import (
	"fmt"
	"net"

	"github.com/abhishek8841/lantern/internal/protocol"
)

func SendMessage(address string, message string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("connect to peer %s: %w", address, err)
	}
	defer conn.Close()

	// _, err = conn.Write([]byte(message))
	// if err != nil {
	// 	return fmt.Errorf("send message: %w", err)
	// }
	if err := protocol.WriteFrame(conn, []byte(message)); err != nil {
		return fmt.Errorf("send framed message: %w", err)
	}

	fmt.Println("Message sent to", address)
	return nil
}
