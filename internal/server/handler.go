package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type input struct {
	URL string `json:"url"`
}

func getURLs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("going to get all shortened urls"))
}

func (s *Server) CreateShortenedURL(w http.ResponseWriter, r *http.Request) {
	newInput := &input{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(newInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	newURL, err := s.URLController.CreateShortenedURL(newInput.URL)
	if err != nil {
		fmt.Printf("error creating shortened URL: %+v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(newURL)
}

func (s *Server) RedirectURL(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "shortenedURL")
	originalURL, err := s.URLController.GetOriginalURL(param)
	if err != nil {
		fmt.Printf("error getting original URL: %+v \n", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func (s *Server) GetUsage(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "shortenedURL")
	usage, err := s.URLController.GetUsage(param)
	if err != nil {
		fmt.Printf("error getting usage: %+v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(usage)
}

func (s *Server) DeleteURL(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "shortenedURL")
	err := s.URLController.DeleteURL(param)
	if err != nil {
		fmt.Printf("error deleting shortened url: %+v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
