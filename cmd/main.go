package main

import (
	"bank-api/internal/http-server/server"
	"bank-api/internal/storage/postgres"
)

func main() {
	dbConn := postgres.New("postgres://slava:slava2000@localhost:5432/bank")
	_ = dbConn

	srv := server.New(dbConn)
	_ = srv.Start()

}
