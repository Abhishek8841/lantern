package main

import (
	"log"

	"github.com/abhishek8841/lantern/internal/transport/tcp"
)

func main() {
	if err := tcp.StartListener(":8080"); err != nil {
		log.Fatal(err)
	}
}
