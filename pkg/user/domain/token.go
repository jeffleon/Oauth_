package domain

import (
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	ConstRefreshTokenKey  = "refreshToken"
	ConstPredeterminedKey = "predetermined"
)

type TokenObj struct {
	Predetermined TokenType
	RefreshToken  TokenType
}

type TokenType struct {
	Secret     string
	Expiration time.Duration
}

type Claims struct {
	Email      string `json:"email"`
	UserID     int64  `json:"user_id"`
	Authorized bool   `json:"authorized"`
	Exp        int64  `json:"exp"`
	jwt.Claims
}

type TokenRepository interface {
	CreateToken(*User, string) (string, error)
	VerifyToken(string, string) (*Claims, error)
}

func (t TokenObj) GetPredeterminedSecret() string {
	return t.Predetermined.Secret
}

func (t TokenObj) GetPredeterminedExpiration() time.Duration {
	return t.Predetermined.Expiration
}

func (t TokenObj) GetRefreshTokenExpiration() time.Duration {
	return t.RefreshToken.Expiration
}

func (t TokenObj) GetRefreshTokenSecret() string {
	return t.RefreshToken.Secret
}
