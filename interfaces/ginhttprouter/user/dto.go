package user

import (
	"github.com/gofrs/uuid"
	"time"
)

type User struct {
	Id           uuid.UUID `json:"id"`
	Username     string    `json:"username" binding:"required"`
	Email        string    `json:"email,omitempty"`
	RoleId       string    `json:"roleId,omitempty"`
	DateCreated  time.Time `json:"dateCreated"`
	DateModified time.Time `json:"dateModified"`
}
