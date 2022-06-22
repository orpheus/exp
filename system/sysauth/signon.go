package sysauth

import (
	"errors"
	"fmt"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/core"
	"golang.org/x/crypto/bcrypt"
)

type SignOnInteractor struct {
	UserRepository        UserRepository
	SkillerRepository     core.SkillerRepository
	PermissionsRepository PermissionRepository
	RoleRepository        RoleRepository
	JWTService            JWTService
	Logger                api.Logger
}

type SignOnRepository interface {
	Login(usernameOrEmail string, password string) (User, error)
	SignUp(user RegisterUser) (User, error)
}

func (s *SignOnInteractor) Login(usernameOrEmail string, password string) (*User, error) {
	skiller, err := s.SkillerRepository.FindByUsername(usernameOrEmail)
	if err != nil {
		skiller, err = s.SkillerRepository.FindByEmail(usernameOrEmail)
		if err != nil {
			return nil, err
		}
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(skiller.Password), []byte(password)); err != nil {
		return nil, errors.New("unauthorized")
	}

	role, err := s.RoleRepository.FindById(skiller.RoleId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Role %s does not exist", role.Id.String()))
	}

	jwt := s.JWTService.GenerateToken(skiller.Id, role.Permissions)

	role.toDTO()
	user := &User{
		Id:       skiller.Id,
		Username: skiller.Username,
		Email:    skiller.Email,
		Role:     role,
		JWT:      jwt,
	}

	return user, nil
}

func (s *SignOnInteractor) SignUp(user RegisterUser) (*User, error) {
	if len(user.Password) == 0 {
		return nil, errors.New("missing password")
	}

	if user.RoleId.IsNil() && user.RoleName == "" {
		return nil, errors.New("must provide either `roleId` or `roleName`")
	}

	var role Role
	var err error

	if !user.RoleId.IsNil() {
		role, err = s.RoleRepository.FindById(user.RoleId)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Role id %s does not exist", user.RoleId))
		}
	}

	if user.RoleName != "" {
		role, err = s.RoleRepository.FindByName(user.RoleName)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Role name %s does not exist", user.RoleName))
		}
		user.RoleId = role.Id
	}

	// The second argument is the cost of hashing, which we arbitrarily set as 8
	// (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	createdUser, err := s.UserRepository.CreateOne(user)
	if err != nil {
		return nil, err
	}

	role.toDTO()
	createdUser.Role = role

	return createdUser, nil
}
