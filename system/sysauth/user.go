package sysauth

import (
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/core"
	"github.com/orpheus/exp/system"
	"time"
)

// User is DTO for client views. It is a dumb representation of a Skiller
// with permissions and roles attached for client logic.
type User struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     Role      `json:"role"`
	JWT      string    `json:"accessToken"`
}

// RegisterUser is the intake-DTO for user/skiller registration. Currently, allow
// user to assign a role via an `RoleId` or `RoleName` as both have unique db constraints.
type RegisterUser struct {
	Username string    `json:"username" binding:"required"`
	Password string    `json:"password"  binding:"required"`
	Email    string    `json:"email" binding:"required"`
	RoleId   uuid.UUID `json:"roleId"`
	RoleName string    `json:"roleName"`
}

// UserInteractor implements the UserRepository for use in signon logic.
type UserInteractor struct {
	SkillerInteractor *system.SkillerInteractor
	Logger            api.Logger
}

// UserRepository should just be used for creating and editing individual users
// which includes system-level bindings. It wraps the SkillerInteractor to create
// the domain-level-user, but handles registration, permissions, roles, etc.
type UserRepository interface {
	CreateOne(user RegisterUser) (*User, error)
}

func (u *UserInteractor) ExistsByUsername(username string) bool {
	_, err := u.SkillerInteractor.FindByUsername(username)
	if err != nil {
		return false
	}
	return true
}

func (u *UserInteractor) CreateOne(user RegisterUser) (*User, error) {
	skiller, err := u.SkillerInteractor.CreateOne(core.Skiller{
		Email:        user.Email,
		Username:     user.Username,
		Password:     user.Password,
		RoleId:       user.RoleId,
		DateCreated:  time.Time{},
		DateModified: time.Time{},
	})
	if err != nil {
		return nil, err
	}
	return &User{
		Id:       skiller.Id,
		Username: skiller.Username,
		Email:    skiller.Email,
	}, nil
}
