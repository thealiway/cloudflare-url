package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type input struct {
	URL string `json:"url"`
}

func getURLs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("going to get all shortened urls"))
}

func (s *Server) CreateShortenedURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in correct handler")
	newInput := &input{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(newInput)
	if err != nil {
		w.WriteHeader(400)
	}

	fmt.Println("about to create")

	newURL, err := s.URLController.CreateShortenedURL(newInput.URL)
	if err != nil {
		fmt.Printf("error creating shortened URL: %+v \n", err)
		w.WriteHeader(500)
	}
	json.NewEncoder(w).Encode(newURL)
}
