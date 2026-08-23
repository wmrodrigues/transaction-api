package authentication

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	// here we could add more data, in this case for simplicity we only have the user id and email
	UserID string `json:"user_id"`
	Email  string `json:"email"`

	jwt.RegisteredClaims
}

type claimsKey struct{}

func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(*Claims)
	return claims, ok
}
