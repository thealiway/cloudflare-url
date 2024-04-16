package server

import (
	"fmt"
	"net/http"
)

func getURLs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("going to get all shortened urls"))
}

func (s *Server) createShortenedURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	newURL := s.URLController.CreateShortenedURL("test url")
	w.Write([]byte(newURL.ShortenedURL))
}
