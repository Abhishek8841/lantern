package main

import (
	"log"
	"os"

	"github.com/abhishek8841/lantern/internal/peer"
	"github.com/abhishek8841/lantern/internal/transport/tcp"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: lantern listen | send")
	}

	switch os.Args[1] {
	case "listen":
		manager := peer.NewManager()
		if err := tcp.StartListener(":8080", manager); err != nil {
			log.Fatal(err)
		}
	case "send":
		peer, err := tcp.DialPeer("127.0.0.1:8080")
		if err != nil {
			log.Fatal(err)
		}
		defer peer.Close()

		messages := []string{
			"Hello from LANtern!",
			"This is message two",
			"This is message three",
		}

		for _, message := range messages {
			if err := peer.SendMessage(message); err != nil {
				log.Fatal(err)
			}
		}
	default:
		log.Fatal("unknown mode: use listen or send")
	}
}
