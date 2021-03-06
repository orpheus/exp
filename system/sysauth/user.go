package sysauth

import (
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/core"
	"github.com/orpheus/exp/system"
	"time"
)

// User is meant to be a dumb representation of the Skiller struct
// used in the business logic for dealing with user-based interactions.
type User struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username" binding:"required"`
	Password string    `json:"password,omitempty"`
	Email    string    `json:"email" binding:"required"`
	RoleId   uuid.UUID `json:"roleId,omitempty" binding:"required"`
	RoleName string    `json:"role,omitempty" binding:"required"`
	JWT      string    `json:"accessToken,omitempty"`
}

type RegisterUser struct {
	Username string    `json:"username" binding:"required"`
	Password string    `json:"password"  binding:"required"`
	Email    string    `json:"email" binding:"required"`
	RoleId   uuid.UUID `json:"roleId"`
	RoleName string    `json:"roleName"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"  binding:"required"`
}

// UserInteractor implements the UserRepository for use in signon logic.
type UserInteractor struct {
	SkillerInteractor *system.SkillerInteractor
	Logger            api.Logger
}

// UserRepository is an overlay repository over Skiller. The actual
// UserInteractor implements all of these  methods by calling its
// SkillerInteractor and transforming into User structs.
//
// The reason to do this for now is to keep User and Skiller separate
// without having them both live in the database. So for now, Skiller
// is the source of truth and the User is the DTO wrapper.
type UserRepository interface {
	FindAll() ([]*User, error)
	FindById(id uuid.UUID) (User, error)
	FindByUsername(username string) (User, error)
	FindByEmail(email string) (User, error)
	CreateOne(user RegisterUser) (User, error)
}

func (u *User) RemovePassword() {
	u.Password = ""
}

func (u *UserInteractor) FindAll() ([]*User, error) {
	skillers, err := u.SkillerInteractor.FindAll()
	if err != nil {
		return nil, err
	}
	var users []*User
	for _, skiller := range skillers {
		users = append(users, &User{
			Id:       skiller.Id,
			Username: skiller.Username,
			Password: skiller.Password,
			Email:    skiller.Email,
			RoleId:   skiller.RoleId,
		})
	}
	return users, nil
}

func (u *UserInteractor) FindById(id uuid.UUID) (User, error) {
	skiller, err := u.SkillerInteractor.FindById(id)
	if err != nil {
		return User{}, err
	}
	return User{
		Id:       skiller.Id,
		Username: skiller.Username,
		Password: skiller.Password,
		Email:    skiller.Email,
		RoleId:   skiller.RoleId,
	}, nil
}

func (u *UserInteractor) FindByUsername(username string) (User, error) {
	skiller, err := u.SkillerInteractor.FindByUsername(username)
	if err != nil {
		return User{}, err
	}
	return User{
		Id:       skiller.Id,
		Username: skiller.Username,
		Password: skiller.Password,
		Email:    skiller.Email,
		RoleId:   skiller.RoleId,
	}, nil
}

func (u *UserInteractor) FindByEmail(email string) (User, error) {
	skiller, err := u.SkillerInteractor.FindByEmail(email)
	if err != nil {
		return User{}, err
	}
	return User{
		Id:       skiller.Id,
		Username: skiller.Username,
		Password: skiller.Password,
		Email:    skiller.Email,
		RoleId:   skiller.RoleId,
	}, nil
}

func (u *UserInteractor) CreateOne(user RegisterUser) (User, error) {
	skiller, err := u.SkillerInteractor.CreateOne(core.Skiller{
		Email:        user.Email,
		Username:     user.Username,
		Password:     user.Password,
		RoleId:       user.RoleId,
		DateCreated:  time.Time{},
		DateModified: time.Time{},
	})
	if err != nil {
		return User{}, err
	}
	return User{
		Id:       skiller.Id,
		Username: skiller.Username,
		Password: "",
		Email:    skiller.Email,
		RoleId:   skiller.RoleId,
	}, nil
}
