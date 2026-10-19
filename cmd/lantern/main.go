package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/abhishek8841/lantern/internal/identity"
	"github.com/abhishek8841/lantern/internal/peer"
	"github.com/abhishek8841/lantern/internal/transport/tcp"
)

func main() {

	listenAddr := flag.String("listen", ":8080", "TCP listen address")
	identityPath := flag.String("identity", "lantern-identity", "Identity file path")
	connectAddr := flag.String("connect", "", "Address of peer to connect to")
	flag.Parse()

	peerID, err := identity.LoadOrCreate(*identityPath)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	manager := peer.NewManager()

	go func() {
		if err := tcp.StartListener(ctx, *listenAddr, manager, peerID); err != nil {
			log.Printf("listener stopped: %v", err)
			stop()
		}
	}()

	if *connectAddr == "" {
		fmt.Println("Listening only. Press Ctrl+C to quit.")
		<-ctx.Done()
		return
	}

	conn, err := tcp.DialPeer(*connectAddr)
	if err != nil {
		log.Fatalf("dial peer: %v", err)
	}

	result := make(chan tcp.ConnectionResult, 1)

	go tcp.ServeConnection(conn, peerID, ctx, manager, result)

	var remoteID string

	select {
	case connection := <-result:
		if connection.Err != nil {
			log.Fatalf("connect to peer: %v", connection.Err)
		}
		remoteID = connection.PeerID
		fmt.Println("Connected to peer:", connection.PeerID)

	case <-ctx.Done():
		return
	}

	fmt.Println("Type a message and press Enter. Use /quit to exit.")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		message := scanner.Text()

		if message == "/quit" {
			break
		}

		if err := manager.SendMessage(peerID, remoteID, message); err != nil {
			log.Printf("send MESSAGE: %v", err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("read input: %v", err)
	}

}
