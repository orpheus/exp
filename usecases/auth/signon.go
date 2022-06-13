package usecases

import (
	"errors"
	"fmt"
	"github.com/orpheus/exp/interfaces"
	"golang.org/x/crypto/bcrypt"
)

type SignOnInteractor struct {
	UserRepository UserRepository
	RoleRepository RoleRepository
	JWTService     JWTService
	Logger         interfaces.Logger
}

func (s *SignOnInteractor) Login(username string, password string) (User, error) {
	user, err := s.UserRepository.FindByUsername(username)

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return User{}, errors.New("unauthorized")
	}

	role, err := s.RoleRepository.FindById(user.RoleId)
	if err != nil {
		return User{}, errors.New(fmt.Sprintf("Role %s does not exist", role.Id.String()))
	}

	jwt := s.JWTService.GenerateToken(user.Id, role.Permissions)
	user.JWT = jwt

	return user, nil
}

func (s *SignOnInteractor) SignUp(user User) (User, error) {
	if len(user.Password) == 0 {
		return User{}, errors.New("missing password")
	}

	_, err := s.RoleRepository.FindById(user.RoleId)
	if err != nil {
		return User{}, errors.New(fmt.Sprintf("Role %s does not exist", user.RoleId))
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
		return User{}, nil
	}

	createdUser.RemovePassword()

	return createdUser, nil
}
