package sysauth

import (
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
)

const admin = "admin"
const adminEmail = "titanroark@gmail.com"

func (r *RoleInteractor) CreateAdminRole() uuid.UUID {
	role, err := r.FindByName(admin)
	if err != nil {
		role, err = r.CreateOne(Role{
			Name:        admin,
			Permissions: []string{"*"},
		})
		if err != nil {
			log.Panicf("Failed to create admin role: %s", err.Error())
		}
		return role.Id
	}
	return role.Id
}

func (u *UserInteractor) CreateAdminUser(adminRoleId uuid.UUID) {
	exists := u.ExistsByUsername(admin)
	if exists == false {
		// The second argument is the cost of hashing, which we arbitrarily set as 8
		// (this value can be more or less, depending on the computing power you wish to utilize)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin), 8)
		if err != nil {
			log.Panicf("Failed to create admin user: %s", err.Error())
		}
		_, err = u.CreateOne(RegisterUser{
			Username: admin,
			Password: string(hashedPassword),
			RoleName: admin,
			RoleId:   adminRoleId,
			Email:    adminEmail,
		})
		if err != nil {
			log.Panicf("Failed to create admin user: %s", err.Error())
		}
	}
}
