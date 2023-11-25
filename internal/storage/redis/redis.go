package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Session struct {
	ID          uuid.UUID `json:"id"`
	SessionArgs `json:"session_args"`
}

type SessionArgs struct {
	Email        string    `json:"email"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type RefreshSession struct {
	Client *redis.Client
}

func New() *RefreshSession {
	conn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	log.Println(conn.Ping(context.Background()))
	return &RefreshSession{
		Client: conn,
	}
}

func sessionIDKey(id uuid.UUID) string {
	return fmt.Sprintf("session: %d", id)
}

func (s *RefreshSession) CreateSession(id uuid.UUID, arg SessionArgs) (Session, error) {
	sessionData := SessionArgs{
		Email:        arg.Email,
		RefreshToken: arg.RefreshToken,
		UserAgent:    arg.UserAgent,
		ClientIp:     arg.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    arg.ExpiresAt,
	}

	data, err := json.Marshal(sessionData)
	log.Println(data)
	if err != nil {
		return Session{}, err
	}

	sessionKey := sessionIDKey(id)

	res := s.Client.SetNX(context.Background(), sessionKey, string(data), 0)
	if err := res.Err(); err != nil {
		return Session{}, err
	}
	return Session{id, sessionData}, nil
}
