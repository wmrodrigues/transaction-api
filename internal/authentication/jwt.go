package authentication

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateToken(userID string, email string) (RegisteredToken, error)
	ValidateToken(token string) (*Claims, error)
}

type JWTService struct {
	secret []byte
}

type RegisteredToken struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	IssuedAt  string `json:"issued_at"`
}

const TokenDateFormat = "2006-01-02 15:04:05"

func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: []byte(secret)}
}

func (s *JWTService) GenerateToken(userID string, email string) (RegisteredToken, error) {
	var registeredToken RegisteredToken
	expiresAt := jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	issuedAt := jwt.NewNumericDate(time.Now())
	claims := Claims{
		UserID:           userID,
		Email:            email,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: expiresAt, IssuedAt: issuedAt},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return registeredToken, fmt.Errorf("error signing token: %w", err)
	}
	registeredToken = RegisteredToken{
		Token:     signedToken,
		ExpiresAt: expiresAt.Time.Format(TokenDateFormat),
		IssuedAt:  issuedAt.Time.Format(TokenDateFormat),
	}
	return registeredToken, nil
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) { return s.secret, nil },
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
