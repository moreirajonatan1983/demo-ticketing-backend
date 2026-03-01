package ports

import "github.com/demoticketing/auth/internal/core/domain"

// UserRepository define el contrato de persistencia para usuarios.
type UserRepository interface {
	Save(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}
