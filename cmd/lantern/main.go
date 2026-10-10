package main

import (
	"log"
	"os"

	"github.com/abhishek8841/lantern/internal/transport/tcp"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: lantern listen | send")
	}

	switch os.Args[1] {
	case "listen":
		if err := tcp.StartListener(":8080"); err != nil {
			log.Fatal(err)
		}
	case "send":
		if err := tcp.SendMessage("127.0.0.1:8080", "Hello from LANtern!"); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("unknown mode: use listen or send")
	}
}
