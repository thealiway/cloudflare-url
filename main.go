package main

import (
	"cloudflareurl/internal/server"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	server, err := server.NewServer()
	if err != nil {
		log.Println("Unable to connect to DB")
		log.Fatal(err)
	}

	err = http.ListenAndServe(":"+port, server)
	if err != nil {
		log.Fatal(err)
	}
}
