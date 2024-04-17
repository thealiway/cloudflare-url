package server

import (
	"cloudflareurl/internal/controllers"
	"cloudflareurl/internal/store/urls"
	"net/http"

	"github.com/go-chi/chi"
)

type Serverer interface {
	CreateShortenedURL(w http.ResponseWriter, r *http.Request)
}

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

	u, err := controllers.NewURLController(URLStore)
	if err != nil {
		return nil, err
	}
	r := chi.NewRouter()

	server := &Server{
		router:        r,
		URLController: u,
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I'm up"))
	})

	r.Get("/shortenedURLs", getURLs)

	r.Post("/shortenedURL", server.CreateShortenedURL)

	r.Route("/s", func(r chi.Router) {
		r.Route("/{shortenedURL}", func(r chi.Router) {
			r.Get("/", server.RedirectURL)
			r.Get("/usage", server.GetUsage)
			r.Delete("/", server.DeleteURL)
		})
	})

	return server, nil
}
