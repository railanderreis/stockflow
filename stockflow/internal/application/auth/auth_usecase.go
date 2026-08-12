package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Interfaces de Repositório e Serviços
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	Save(ctx context.Context, user *User) error
}

type User struct {
	ID           string   `json:"id"`
	Email        string   `json:"email"`
	PasswordHash string   `json:"-"`
	Name         string   `json:"name"`
	Role         string   `json:"role"`
	Permissions  []string `json:"permissions"`
}

// DTOs (Data Transfer Objects)
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterInput struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type AuthOutput struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresIn    int64   `json:"expires_in"`
	User         UserDTO `json:"user"`
}

type UserDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type Claims struct {
	UserID      string   `json:"user_id"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// Casos de Uso
type AuthUseCase struct {
	userRepo  UserRepository
	jwtSecret []byte
}

func NewAuthUseCase(userRepo UserRepository, jwtSecret string) *AuthUseCase {
	return &AuthUseCase{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

// Execute Login
func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (*AuthOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil || user == nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Compara hash da senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	return uc.generateTokenPair(user)
}

// Execute Register
func (uc *AuthUseCase) Register(ctx context.Context, input RegisterInput) (*UserDTO, error) {
	existing, _ := uc.userRepo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return nil, errors.New("e-mail já cadastrado")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("erro ao processar senha: %w", err)
	}

	newUser := &User{
		ID:           fmt.Sprintf("usr_%d", time.Now().UnixNano()), // Exemplo de ID gerado
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Name:         input.Name,
		Role:         input.Role,
		Permissions:  input.Permissions,
	}

	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		return nil, fmt.Errorf("erro ao salvar usuário: %w", err)
	}

	return &UserDTO{
		ID:          newUser.ID,
		Name:        newUser.Name,
		Email:       newUser.Email,
		Role:        newUser.Role,
		Permissions: newUser.Permissions,
	}, nil
}

// Execute Refresh Token
func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*AuthOutput, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return uc.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("refresh token inválido ou expirado")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("claims inválidas no token")
	}

	user, err := uc.userRepo.FindByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return nil, errors.New("usuário não encontrado")
	}

	return uc.generateTokenPair(user)
}

// Função auxiliar interna para emissão de JWT
func (uc *AuthUseCase) generateTokenPair(user *User) (*AuthOutput, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:      user.ID,
		Email:       user.Email,
		Role:        user.Role,
		Permissions: user.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(uc.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Refresh token válido por 7 dias
	refreshClaims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(uc.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    expirationTime.Unix(),
		User: UserDTO{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			Role:        user.Role,
			Permissions: user.Permissions,
		},
	}, nil
}
