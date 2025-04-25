package main

import (
	"log"
	"personal_blog/server"
)

func main() {
	srv, err := server.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Start(); err != nil {
		log.Fatal("Server error:", err)
	}
	log.Fatal("good")
}
