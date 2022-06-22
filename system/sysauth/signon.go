package sysauth

import (
	"errors"
	"fmt"
	"github.com/orpheus/exp/api"
	"golang.org/x/crypto/bcrypt"
)

type SignOnInteractor struct {
	UserRepository UserRepository
	RoleRepository RoleRepository
	JWTService     JWTService
	Logger         api.Logger
}

type SignOnRepository interface {
	Login(username string, password string) (User, error)
	SignUp(user User) (User, error)
}

func (s *SignOnInteractor) Login(usernameOrEmail string, password string) (User, error) {
	user, err := s.UserRepository.FindByUsername(usernameOrEmail)
	if err != nil {
		user, err = s.UserRepository.FindByEmail(usernameOrEmail)
		if err != nil {
			return User{}, err
		}
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return User{}, errors.New("unauthorized")
	}

	// After we've verified the password, remove it, so it doesn't get sent over the wire
	user.RemovePassword()

	role, err := s.RoleRepository.FindById(user.RoleId)
	if err != nil {
		return User{}, errors.New(fmt.Sprintf("Role %s does not exist", role.Id.String()))
	}

	jwt := s.JWTService.GenerateToken(user.Id, role.Permissions)
	user.JWT = jwt

	return user, nil
}

func (s *SignOnInteractor) SignUp(user RegisterUser) (User, error) {
	if len(user.Password) == 0 {
		return User{}, errors.New("missing password")
	}

	if user.RoleId.IsNil() && user.RoleName == "" {
		return User{}, errors.New("must provide either `roleId` or `roleName`")
	}

	if !user.RoleId.IsNil() {
		_, err := s.RoleRepository.FindById(user.RoleId)
		if err != nil {
			return User{}, errors.New(fmt.Sprintf("Role id %s does not exist", user.RoleId))
		}
	}

	if user.RoleName != "" {
		role, err := s.RoleRepository.FindByName(user.RoleName)
		if err != nil {
			return User{}, errors.New(fmt.Sprintf("Role name %s does not exist", user.RoleName))
		}
		user.RoleId = role.Id
	}

	// The second argument is the cost of hashing, which we arbitrarily set as 8
	// (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return User{}, err
	}

	user.Password = string(hashedPassword)
	createdUser, err := s.UserRepository.CreateOne(user)
	if err != nil {
		return User{}, err
	}

	createdUser.RemovePassword()

	return createdUser, nil
}
