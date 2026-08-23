package authentication

import (
	"fmt"
	"net/http"
	"strings"
	"transaction-api/internal/common"
)

func AuthMiddleware(tokenService TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					common.SendUnauthorizedResponse(w, fmt.Errorf("missing authorization header"))
					return
				}
				tokenString := strings.TrimPrefix(authHeader, "Bearer ")
				claims, err := tokenService.ValidateToken(tokenString)
				if err != nil {
					common.SendUnauthorizedResponse(w, fmt.Errorf("invalid token: %w", err))
					return
				}
				ctx := ContextWithClaims(r.Context(), claims)
				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}
