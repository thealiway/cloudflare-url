package server

import (
	"cloudflareurl/internal/controllers"
	"cloudflareurl/internal/store/urls"
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	router          *chi.Mux
	URLController   *controllers.URLController
	UsageController *controllers.UsageController
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func NewServer() (*Server, error) {
	store, err := urls.NewStore()
	if err != nil {
		return nil, err
	}

	u := controllers.NewURLController(store)
	a := controllers.NewUsageController(store)

	r := chi.NewRouter()

	server := &Server{
		router:          r,
		URLController:   u,
		UsageController: a,
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I'm up"))
	})

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
