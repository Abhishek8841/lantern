package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

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

	// if *connectAddr == "" {
	// 	fmt.Println("Listening only. Press Ctrl+C to quit.")
	// 	<-ctx.Done()
	// 	return
	// }

	var remoteID string
	if *connectAddr != "" {
		conn, err := tcp.DialPeer(*connectAddr)
		if err != nil {
			log.Fatalf("dial peer: %v", err)
		}

		result := make(chan tcp.ConnectionResult, 1)

		go tcp.ServeConnection(conn, peerID, ctx, manager, result)

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
	}

	fmt.Println("Type a message and press Enter. Use /quit to exit.")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		// fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		message := scanner.Text()

		if message == "/quit" {
			break
		}

		if message == "/exit-peer-chat" {
			remoteID = ""
			continue
		}

		if message == "/peers" {
			ids := manager.List()

			if len(ids) == 0 {
				fmt.Println("No peers connected.")
				continue
			}

			fmt.Println("Connected peers:")
			for _, id := range ids {
				fmt.Println(" -", id)
			}
			continue
		}

		if len(message) > 6 && message[:6] == "/chat " {
			targetID := message[6:]

			if _, ok := manager.Get(targetID); !ok {
				fmt.Println("Peer not found. Use /peers to see connected peers.")
				continue
			}

			remoteID = targetID
			fmt.Println("Now chatting with:", remoteID)
			continue
		}

		if strings.HasPrefix(message, "/connect ") {
			addr := strings.TrimSpace(strings.TrimPrefix(message, "/connect "))
			if addr == "" {
				fmt.Println("Usage: /connect <host:port>")
				continue
			}
			go func(addr string) {
				conn, err := tcp.DialPeer(addr)
				if err != nil {
					fmt.Printf("dial peer: %v", err)
					return
				}

				result := make(chan tcp.ConnectionResult, 1)

				go tcp.ServeConnection(conn, peerID, ctx, manager, result)

				select {
				case connection := <-result:
					if connection.Err != nil {
						fmt.Printf("connect to peer: %v", connection.Err)
						return
					}
					fmt.Println("Connected to peer:", connection.PeerID)

				case <-ctx.Done():
					return
				}
			}(addr)
			continue
		}

		if remoteID != "" {
			if err := manager.SendMessage(peerID, remoteID, message); err != nil {
				log.Printf("send MESSAGE: %v", err)
				remoteID = ""
				continue
			}
		} else {
			fmt.Println("No peer selected. Use /peers, then /chat <peer-id>.")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("read input: %v", err)
	}

}
