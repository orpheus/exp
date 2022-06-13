package core

import (
	"github.com/gofrs/uuid"
	"time"
)

type SkillerRepository interface {
	FindAll() ([]Skiller, error)
	FindById(id uuid.UUID) (Skiller, error)
	FindByUsername(username string) (Skiller, error)
	FindByEmail(email string) (Skiller, error)
	CreateOne(skiller Skiller) (Skiller, error)
}

// Skiller is the entity which acts in the system and has skills.
type Skiller struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email" binding:"required"`
	Username string    `json:"username" binding:"required"`
	Password string    `json:"password,omitempty"`
	// Here until I understand how to separate a `Skiller` and a `User`
	RoleId       uuid.UUID `json:"roleId,omitempty" binding:"required"`
	DateCreated  time.Time `json:"dateCreated"`
	DateModified time.Time `json:"dateModified"`
}
