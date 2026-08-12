package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/railanderreis/stockflow/stockflow/internal/application/auth"
)

type AuthMiddleware struct {
	jwtSecret []byte
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: []byte(secret),
	}
}

// Authenticate valida se o token Bearer no Header Authorization é válido
func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "cabeçalho de autorização ausente"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "formato de token inválido"})
		}

		tokenString := parts[1]
		claims := &auth.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return m.jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token inválido ou expirado"})
		}

		// Armazena no contexto local da requisição HTTP
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("permissions", claims.Permissions)

		return c.Next()
	}
}

// RequirePermission verifica se o usuário autenticado possui a permissão necessária
func (m *AuthMiddleware) RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permissionsRaw := c.Locals("permissions")
		if permissionsRaw == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "acesso negado"})
		}

		permissions, ok := permissionsRaw.([]string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "permissões inválidas"})
		}

		hasPermission := false
		for _, p := range permissions {
			if p == permission || p == "admin:all" {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "permissão insuficiente para este recurso"})
		}

		return c.Next()
	}
}
