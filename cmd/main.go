package main

import (
	"bank-api/internal/config"
	"bank-api/internal/http-server/server"
	"bank-api/internal/storage/postgres"
	"bank-api/internal/storage/redis"
	"log"
)

func main() {
	cfg := config.MustLoad()

	dbConn := postgres.New(cfg.PostgresURL)
	sConn := redis.New()

	srv := server.New(dbConn, sConn, cfg)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

}
