package server

import (
	"github.com/go-chi/chi"
)

type Server struct{}

func NewServer() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/shortenedURLs", getURLs)

	r.Post("/shortenedURL", createShortenedURL)

	return r
}
