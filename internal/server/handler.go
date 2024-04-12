package server

import (
	"net/http"
)

func getURLs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("going to get all shortened urls"))
}

func createShortenedURL(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("going to create a shortened url and return to user"))
}
