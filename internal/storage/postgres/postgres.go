package postgres

import (
	"bank-api/internal/model"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Storage struct {
	Conn *pgxpool.Pool
}

func New(dbUrl string) *Storage {
	poolCfg, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatalln("ERROR: parse db url")
	}

	connPool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		log.Fatalln("ERROR: connect to db ")
	}
	//conn, err := pgx.Connect(context.Background(), dbUrl)

	return &Storage{Conn: connPool}
}

func (s *Storage) UserExists(email string) (bool, error) {
	var count int

	sqlString := "SELECT count(*) FROM public.users where email = $1"
	row := s.Conn.QueryRow(context.Background(), sqlString, email)
	if _ = row.Scan(&count); count > 0 {
		return false, errors.New("user with this email exists")
	}
	//if count > 0 {
	//	return false, errors.New("user with this email exists")
	//}

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
	var user model.User
	sqlString := "SELECT email, encrypted_password, created_at, full_name FROM public.users WHERE email = $1"

	row := s.Conn.QueryRow(context.Background(), sqlString, email)
	err := row.Scan(&user.Email, &user.EncryptedPassword, &user.CreatedAt, &user.FullName)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
