package identity

import (
	"context"
	"fmt"
	"time"

	"github.com/railanderreis/stockflow/stockflow/internal/domain/identity"
	"github.com/railanderreis/stockflow/stockflow/internal/infrastructure/auth"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserDTO   `json:"user"`
}

type UserDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginUseCase struct {
	repo       identity.UserRepository
	jwtManager *auth.JWTManager
}

func NewLoginUseCase(repo identity.UserRepository, jwtManager *auth.JWTManager) *LoginUseCase {
	return &LoginUseCase{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	user, err := uc.repo.GetByEmail(input.Email)
	if err != nil {
		return nil, identity.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, identity.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, identity.ErrInvalidCredentials
	}

	permissions, err := uc.repo.GetPermissionsByRoleID(user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user permissions: %w", err)
	}

	token, expiresAt, err := uc.jwtManager.GenerateToken(user.ID, user.Email, user.RoleName, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginOutput{
		Token:     token,
		ExpiresAt: expiresAt,
		User: UserDTO{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.RoleName,
		},
	}, nil
}
