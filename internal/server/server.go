package server

import (
	"cloudflareurl/internal/controllers"
	"cloudflareurl/internal/store/urls"
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	router        *chi.Mux
	URLController *controllers.URLController
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func NewServer() (*Server, error) {
	URLStore, err := urls.NewURLStore()
	if err != nil {
		return nil, err
	}

	u := controllers.NewURLController(URLStore)
	r := chi.NewRouter()

	server := &Server{
		router:        r,
		URLController: u,
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I'm up"))
	})

	r.Get("/shortenedURLs", getURLs)

	r.Post("/shortenedURL", server.createShortenedURL)

	return server, nil
}
