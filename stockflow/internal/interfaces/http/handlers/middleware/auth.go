package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/railanderreis/stockflow/stockflow/internal/infrastructure/auth"
)

type contextKey string

const (
	UserContextKey contextKey = "user_claims"
)

func AuthMiddleware(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Cabeçalho de autorização ausente")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondJSONError(w, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Formato de token inválido. Use 'Bearer <token>'")
				return
			}

			claims, err := jwtManager.ValidateToken(parts[1])
			if err != nil {
				respondJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Token inválido ou expirado")
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
			if !ok {
				respondJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Acesso não autorizado")
				return
			}

			// ADMIN possui acesso irrestrito
			if claims.Role == "ADMIN" {
				next.ServeHTTP(w, r)
				return
			}

			hasPerm := false
			for _, perm := range claims.Permissions {
				if perm == permission {
					hasPerm = true
					break
				}
			}

			if !hasPerm {
				respondJSONError(w, http.StatusForbidden, "FORBIDDEN", "Permissão insuficiente para acessar este recurso")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func respondJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"details": map[string]interface{}{},
		},
	})
}
