package auth_test

import (
	"testing"
	"time"

	"github.com/railanderreis/stockflow/stockflow/internal/infrastructure/auth"
)

func TestJWTManager_GenerateAndValidateToken(t *testing.T) {
	secret := "super-secret-key-stockflow"
	manager := auth.NewJWTManager(secret, 1*time.Hour)

	userID := "usr-123"
	email := "comprador@stockflow.com"
	role := "BUYER"
	permissions := []string{"purchase:read", "purchase:create"}

	token, exp, err := manager.GenerateToken(userID, email, role, permissions)
	if err != nil {
		t.Fatalf("esperado erro nil ao gerar token, obteve: %v", err)
	}

	if token == "" {
		t.Fatal("esperado token válido, obteve string vazia")
	}

	if exp.Before(time.Now()) {
		t.Fatal("data de expiração do token deve ser no futuro")
	}

	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("esperado erro nil ao validar token, obteve: %v", err)
	}

	if claims.UserID != userID || claims.Email != email || claims.Role != role {
		t.Fatalf("claims inconsistentes: %+v", claims)
	}
}
