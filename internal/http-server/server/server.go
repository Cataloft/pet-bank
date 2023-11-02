package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Server struct {
	router *chi.Mux
}

func New() *Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	return &Server{
		router: router,
	}
}
func (s *Server) initHandlers() {
	s.router.HandleFunc("/", s.router.NotFoundHandler())
}

func (s *Server) Start() error {
	s.initHandlers()
	return http.ListenAndServe(":8080", s.router)
}
