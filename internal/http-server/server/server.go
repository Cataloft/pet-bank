package server

import (
	"bank-api/internal/http-server/handlers/register_user"
	"bank-api/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type Server struct {
	router  *chi.Mux
	storage *postgres.Storage
}

func New(storage *postgres.Storage) *Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	return &Server{
		router:  router,
		storage: storage,
	}
}
func (s *Server) initHandlers() {
	s.router.HandleFunc("/register_user", register_user.RegisterUser(s.storage))

}

func (s *Server) Start() error {
	s.initHandlers()
	log.Println("server started on port :8080")
	return http.ListenAndServe(":8080", s.router)
}
