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

func (s *Storage) SaveUser(email string, password string, encryptedPassword []byte) error {
	sqlString := "INSERT INTO public.user (email, password, encrypted_password) VALUES ($1, $2, $3)"

	_, err := s.Conn.Exec(context.Background(), sqlString, email, password, encryptedPassword)
	if err != nil {
		return err
	}

	return nil
}
