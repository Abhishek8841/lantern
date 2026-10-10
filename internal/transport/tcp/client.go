package tcp

import (
	"fmt"
	"net"
)

func SendMessage(address string, message string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("connect to peer %s: %w", address, err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	fmt.Println("Message sent to", address)
	return nil
}
