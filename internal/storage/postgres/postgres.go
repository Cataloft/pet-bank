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

	sqlString := "SELECT count(*) FROM public.users where email = $1"
	row := s.Conn.QueryRow(context.Background(), sqlString, email)
	row.Scan(&count)

	if count > 0 {
		return false, errors.New("user with this email exists")
	}
	return true, nil
}

func (s *Storage) SaveUser(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if exist, err := s.UserExists(user.Email); !exist && err != nil {
		return err
	}

	if err := user.BeforeCreate(); err != nil {
		return err
	}

	sqlString := "INSERT INTO public.users (email, encrypted_password, full_name) VALUES ($1, $2, $3)"
	_, err := s.Conn.Exec(context.Background(), sqlString, user.Email, user.EncryptedPassword, user.FullName)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetUser(email string) (model.User, error) {
	var u model.User
	sqlString := "SELECT email, encrypted_password, created_at, full_name FROM public.users WHERE email = $1"

	row := s.Conn.QueryRow(context.Background(), sqlString, email)
	err := row.Scan(&u.Email, &u.EncryptedPassword, &u.CreatedAt, &u.FullName)

	return u, err
}
