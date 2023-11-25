package login_user

import (
	"bank-api/internal/config"
	"bank-api/internal/model"
	"bank-api/internal/storage/redis"
	"bank-api/pkg/token"
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type Request struct {
	Email    string
	Password string
}

type Response struct {
	SessionID             uuid.UUID
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	User                  model.User
}

type UserGetter interface {
	GetUser(email string) (model.User, error)
}

func LoginUser(userGetter UserGetter, refSession *redis.RefreshSession, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Println("failed to decode request body")
			return
		}

		user, err := userGetter.GetUser(req.Email)
		if err != nil {
			log.Println("failed to get user")
			return
		}

		if !user.ComparePassword(req.Password) {
			log.Println("failed to compare password")
			return
		}

		newToken, err := token.New(cfg.Secret)
		if err != nil {
			log.Println("invalid secret")
			return
		}

		accessToken, accessPayload, err := newToken.CreateToken(
			user.Email,
			user.Role,
			cfg.AccessTokenDuration,
		)
		if err != nil {
			log.Println("failed to create access token")
			return
		}

		refreshToken, refreshPayload, err := newToken.CreateToken(
			user.Email,
			user.Role,
			cfg.AccessTokenDuration,
		)
		if err != nil {
			log.Println("failed to create refresh token")
			return
		}

		session, err := refSession.CreateSession(refreshPayload.ID, redis.SessionArgs{
			Email:        user.Email,
			RefreshToken: refreshToken,
			UserAgent:    r.UserAgent(),
			ClientIp:     r.RemoteAddr,
			IsBlocked:    false,
			ExpiresAt:    refreshPayload.ExpiredAt,
		})
		if err != nil {
			log.Println("failed to create session")
			return
		}

		rsp := Response{
			SessionID:             session.ID,
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  accessPayload.ExpiredAt,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
			User:                  user,
		}

		w.WriteHeader(http.StatusOK)
		resp, _ := json.Marshal(rsp)
		w.Write(resp)

	}
}
