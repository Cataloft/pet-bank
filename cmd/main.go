package main

import (
	"bank-api/internal/config"
	"bank-api/internal/http-server/server"
	"bank-api/internal/storage/postgres"
	"bank-api/internal/storage/redis"
)

func main() {
	cfg := config.MustLoad()

	dbConn := postgres.New(cfg.PostgresURL)
	_ = dbConn

	sConn := redis.New()

	srv := server.New(dbConn, sConn, cfg)
	_ = srv.Start()
}
