package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func New(secret string) (*JWTMaker, error) {
	log.Println(secret)
	if len(secret) < minSecretKeySize {
		return nil, errors.New("invalid key size")
	}
	return &JWTMaker{secret}, nil
}

func (maker *JWTMaker) CreateToken(email string, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(email, role, duration)
	if err != nil {
		return "", nil, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))

	return token, payload, nil
}

//func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
//
//}
