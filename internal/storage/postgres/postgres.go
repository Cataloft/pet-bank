package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	Conn *pgx.Conn
}

func New(dbUrl string) *Storage {
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return nil
	}

	return &Storage{Conn: conn}
}
