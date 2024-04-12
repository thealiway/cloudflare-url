package main

import (
	"cloudflareurl/internal/server"
	"log"
	"net/http"
)

func main() {
	server := server.NewServer()
	err := http.ListenAndServe(":8000", server)
	if err != nil {
		log.Fatal(err)
	}
}
