package tcp

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/abhishek8841/lantern/internal/protocol"
)

func TestTCPFramedMessage(t *testing.T) {
	// :0 tells the OS to automatically chose the port 
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen() error = %v", err)
	}
	defer listener.Close()

	received := make(chan []byte, 1)
	serverErr := make(chan error, 1)

	go func() {

		conn, err := listener.Accept()
		if err != nil {
			serverErr <- fmt.Errorf("listener.Accept(): %w", err)
			return
		}
		defer conn.Close()

		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
		message, err := protocol.ReadFrame(conn)
		if err != nil {
			serverErr <- fmt.Errorf("ReadFrame(): %w", err)
			return
		}

		received <- message
	}()

	conn, err := net.DialTimeout("tcp", listener.Addr().String(), 3*time.Second)
	if err != nil {
		t.Fatalf("net.dialTimeout() error = %v", err)
	}
	defer conn.Close()

	want := []byte("Hello from LANtern!")

	_ = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))

	if err := protocol.WriteFrame(conn, want); err != nil {
		t.Fatalf("WriteFrame() error = %v", err)
	}

	select {
	case got := <-received:
		if string(got) != string(want) {
			t.Fatalf("got %q, want %q", got, want)
		}
	case err := <-serverErr:
		t.Fatalf("server failed: %v", err)
	case <-time.After(4 * time.Second):
		t.Fatal("timed out waiting for server")
	}

}
