package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/demoticketing/auth/internal/core/domain"
	"github.com/demoticketing/auth/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

// AuthService implementa la lógica de negocio de autenticación.
type AuthService struct {
	repo ports.UserRepository
}

func NewAuthService(repo ports.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// RegisterInput contiene los datos del formulario de registro.
type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

// LoginInput contiene las credenciales de inicio de sesión.
type LoginInput struct {
	Email    string
	Password string
}

// AuthResult es la respuesta tras un login/registro exitoso.
type AuthResult struct {
	Token string     `json:"token"`
	User  PublicUser `json:"user"`
}

// PublicUser es la representación pública del usuario (sin password hash).
type PublicUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	ErrEmailAlreadyExists = errors.New("el email ya está registrado")
	ErrInvalidCredentials = errors.New("credenciales incorrectas")
	ErrWeakPassword       = errors.New("la contraseña debe tener al menos 6 caracteres")
)

// Register crea un nuevo usuario con la contraseña hasheada con bcrypt.
func (s *AuthService) Register(input RegisterInput) (*AuthResult, error) {
	if len(input.Password) < 6 {
		return nil, ErrWeakPassword
	}

	// Verificar si el email ya existe
	existing, _ := s.repo.FindByEmail(input.Email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash de la contraseña con bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error al hashear contraseña: %w", err)
	}

	user := &domain.User{
		ID:           fmt.Sprintf("usr_%d", time.Now().UnixMilli()),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	if err := s.repo.Save(user); err != nil {
		return nil, fmt.Errorf("error al guardar usuario: %w", err)
	}

	token := generateToken(user)
	return &AuthResult{
		Token: token,
		User:  PublicUser{ID: user.ID, Name: user.Name, Email: user.Email},
	}, nil
}

// Login autentica al usuario y retorna un token si las credenciales son válidas.
func (s *AuthService) Login(input LoginInput) (*AuthResult, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token := generateToken(user)
	return &AuthResult{
		Token: token,
		User:  PublicUser{ID: user.ID, Name: user.Name, Email: user.Email},
	}, nil
}

// ForgotPassword simula el envío de un email de recuperación.
// Siempre retorna true por seguridad (no revela si el email existe).
func (s *AuthService) ForgotPassword(email string) bool {
	user, _ := s.repo.FindByEmail(email)
	if user != nil {
		// En producción: generar token de recuperación y enviar email via SES
		fmt.Printf("[AUTH] Enviando email de recuperación a: %s\n", email)
	}
	return true
}

// generateToken crea un token simple de desarrollo.
// En producción reemplazar con JWT firmado con RS256.
func generateToken(user *domain.User) string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "ticketera-dev-secret"
	}
	// Token mock: en producción usar github.com/golang-jwt/jwt
	return fmt.Sprintf("mock-%s-%d-%s", user.ID, time.Now().Unix(), secret[:8])
}
