package main

import (
	"go-host/internal/server"
	"log"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatalf("could not start the server: %v", err)
	}
}
