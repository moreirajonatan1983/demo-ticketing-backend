package domain

import "time"

// User representa un usuario registrado en el sistema.
type User struct {
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"` // Nunca serializar el hash en respuestas
	CreatedAt    time.Time `json:"created_at"`
	ID           string    `json:"id"`
}
