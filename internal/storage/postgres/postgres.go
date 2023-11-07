package postgres

import (
	"bank-api/internal/model"
	"context"
	"errors"
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

func (s *Storage) UserExists(email string) (bool, error) {
	var count int

	sqlString := "SELECT count(*) FROM public.user where email = $1"
	_ = s.Conn.QueryRow(context.Background(), sqlString, email).Scan(&count)
	if count > 0 {
		return false, errors.New("user with this email is already registered")
	}
	return true, nil
}

func (s *Storage) SaveUser(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if _, err := s.UserExists(user.Email); err != nil {
		return err
	}

	if err := user.BeforeCreate(); err != nil {
		return err
	}

	sqlString := "INSERT INTO public.user (email, encrypted_password) VALUES ($1, $2)"
	_, err := s.Conn.Exec(context.Background(), sqlString, user.Email, user.EncryptedPassword)
	if err != nil {
		return err
	}

	return nil
}
