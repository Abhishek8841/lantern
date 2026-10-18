package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/abhishek8841/lantern/internal/peer"
	"github.com/abhishek8841/lantern/internal/transport/tcp"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: lantern listen | send")
	}

	switch os.Args[1] {
	case "listen":
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()
		manager := peer.NewManager()
		if err := tcp.StartListener(ctx, ":8080", manager); err != nil {
			log.Fatal(err)
		}
	case "send":
		peerConn, err := tcp.DialPeer("127.0.0.1:8080")
		if err != nil {
			log.Fatal(err)
		}
		defer peerConn.Close()

		fmt.Println("Connected to peer 127.0.0.1:8080")
		fmt.Println("Type a message and press Enter. Use /quit to exit.")

		scanner := bufio.NewScanner(os.Stdin)

		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				return
			}

			message := scanner.Text()

			if message == "/quit" {
				break
			}

			if err := peerConn.SendMessage(message); err != nil {
				log.Printf("send message: %v", err)
				break
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("read input: %v", err)
		}

	default:
		log.Fatal("unknown mode: use listen or send")
	}
}
