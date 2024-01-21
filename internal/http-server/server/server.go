package server

import (
	"bank-api/internal/config"
	"bank-api/internal/http-server/handlers/user/login"
	"bank-api/internal/http-server/handlers/user/register"
	"bank-api/internal/storage/postgres"
	"bank-api/internal/storage/redis"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type Server struct {
	router  *chi.Mux
	storage *postgres.Storage
	session *redis.RedisClient
	config  *config.Config
}

func New(storage *postgres.Storage, session *redis.RedisClient, cfg *config.Config) *Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	return &Server{
		router:  router,
		storage: storage,
		session: session,
		config:  cfg,
	}
}
func (s *Server) initHandlers() {
	s.router.Post("/auth/register", register.Register(s.storage))
	s.router.Post("/auth/login", login.Login(s.storage, s.session, s.config))

}

func (s *Server) Start() error {
	s.initHandlers()
	log.Println("server started on port :8080")
	return http.ListenAndServe(":8080", s.router)
}
