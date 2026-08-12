package auth

import (
	"github.com/gofiber/fiber/v2"
	appAuth "github.com/railanderreis/stockflow/stockflow/internal/application/auth"
)

type AuthHandler struct {
	authUseCase *appAuth.AuthUseCase
}

func NewAuthHandler(authUseCase *appAuth.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input appAuth.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "dados da requisição inválidos"})
	}

	output, err := h.authUseCase.Login(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(output)
}

// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input appAuth.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "dados da requisição inválidos"})
	}

	output, err := h.authUseCase.Register(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(output)
}

// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&input); err != nil || input.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "refresh_token é obrigatório"})
	}

	output, err := h.authUseCase.RefreshToken(c.Context(), input.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(output)
}
