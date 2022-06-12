package auth

import (
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/domain"
	"github.com/orpheus/exp/usecases"
	"time"
)

// UserInteractor implements the UserRepository for use in signon logic.
type UserInteractor struct {
	SkillerInteractor usecases.SkillerInteractor
}

type UserRepository interface {
	FindById(id uuid.UUID) (User, error)
	FindByUsername(username string) (User, error)
	FindByEmail(email string) (User, error)
	CreateOne(user User) (User, error)
}

// User is meant to be a dumb representation of the Skiller struct
// used in the business logic for dealing with user-based interactions.
type User struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username" binding:"required"`
	Password string    `json:"password,omitempty"`
	Email    string    `json:"email" binding:"required"`
	RoleId   uuid.UUID `json:"roleId,omitempty" binding:"required"`
	JWT      string    `json:"accessToken"`
}

func (u *User) RemovePassword() {
	u.Password = ""
}

func (u *UserInteractor) FindById(id uuid.UUID) (User, error) {
	skiller, err := u.SkillerInteractor.FindById(id)
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

func (u *UserInteractor) FindByUsername(username string) (User, error) {
	skiller, err := u.SkillerInteractor.FindByUsername(username)
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

func (u *UserInteractor) FindByEmail(email string) (User, error) {
	skiller, err := u.SkillerInteractor.FindByEmail(email)
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

func (u *UserInteractor) CreateOne(user User) (User, error) {
	skiller, err := u.SkillerInteractor.CreateOne(domain.Skiller{
		Id:           user.Id,
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
