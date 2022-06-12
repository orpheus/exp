package auth

import (
	"fmt"
	"github.com/orpheus/exp/interfaces"
	"github.com/orpheus/exp/interfaces/ginhttprouter/auth"
	"time"
)

type PermissionRepository interface {
	FindAll() ([]Permission, error)
	FindById(id string) (Permission, error)
	CreateOne(id string) (Permission, error)
	DeleteById(id string) error
}

type Permission struct {
	Id          string    `json:"id"`
	DateCreated time.Time `json:"dateCreated"`
}

type PermissionInteractor struct {
	PermissionRepository PermissionRepository
	Logger               interfaces.Logger
}

func (p *PermissionInteractor) FindAll() ([]Permission, error) {
	return p.PermissionRepository.FindAll()
}

func (p *PermissionInteractor) FindById(id string) (Permission, error) {
	return p.PermissionRepository.FindById(id)
}
func (p *PermissionInteractor) CreateOne(id string) (Permission, error) {
	return p.PermissionRepository.CreateOne(id)
}
func (p *PermissionInteractor) DeleteById(id string) error {
	return p.PermissionRepository.DeleteById(id)
}

// EnforcePermissions creates the default/necessary system permissions
func (p *PermissionInteractor) EnforcePermissions(guardian auth.PermissionGuardian) {
	allPermissions := guardian.GetPermissions()
	existingPermissions, err := p.FindAll()
	if err != nil {
		panic(err)
	}

	_ = p.Logger.Log(fmt.Sprintf("%s", existingPermissions))

	mappedPermissions := make(map[string]bool)
	for _, v := range existingPermissions {
		mappedPermissions[v.Id] = true
	}
	for _, permissionId := range allPermissions {
		if !mappedPermissions[permissionId] {
			_, err := p.CreateOne(permissionId)
			if err != nil {
				panic(err)
			}
		}
	}
}
