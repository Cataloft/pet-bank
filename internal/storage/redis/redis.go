package redis

import (
	"context"
	"encoding/json"
	"errors"
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

type RedisClient struct {
	Client *redis.Client
}

func New() *RedisClient {
	conn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	log.Println(conn.Ping(context.Background()))
	return &RedisClient{
		Client: conn,
	}
}

func sessionIDKey(id uuid.UUID) string {
	return fmt.Sprintf("session: %d", id)
}

func (s *RedisClient) CreateSession(ctx context.Context, id uuid.UUID, arg SessionArgs) (Session, error) {
	data, err := json.Marshal(arg)
	log.Println(data)
	if err != nil {
		return Session{}, err
	}

	sessionKey := sessionIDKey(id)

	res := s.Client.SetNX(ctx, sessionKey, string(data), 0)
	if err := res.Err(); err != nil {
		return Session{}, err
	}
	return Session{id, arg}, nil
}

func (s *RedisClient) GetSessionByID(ctx context.Context, id uuid.UUID) (Session, error) {
	sessionKey := sessionIDKey(id)

	data, err := s.Client.Get(ctx, sessionKey).Result()
	if errors.Is(err, redis.Nil) {
		return Session{}, errors.New("session does not exist")
	} else if err != nil {
		return Session{}, err
	}

	var sessionData SessionArgs
	err = json.Unmarshal([]byte(data), &sessionData)
	if err != nil {
		return Session{}, err
	}

	return Session{id, sessionData}, nil
}

func (s *RedisClient) DeleteSessionByID(ctx context.Context, id uuid.UUID) error {
	sessionKey := sessionIDKey(id)

	err := s.Client.Del(ctx, sessionKey).Err()
	if errors.Is(err, redis.Nil) {
		return errors.New("session does not exist")
	} else if err != nil {
		return err
	}

	return nil
}

func (s *RedisClient) UpdateSession(ctx context.Context, id uuid.UUID, arg SessionArgs) error {
	data, err := json.Marshal(arg)
	if err != nil {
		return err
	}

	sessionKey := sessionIDKey(id)

	err = s.Client.SetXX(ctx, sessionKey, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return errors.New("session does not exist")
	} else if err != nil {
		return err
	}

	return nil
}
