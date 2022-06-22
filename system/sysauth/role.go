package sysauth

import (
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/api"
	"time"
)

// RoleInteractor Service
type RoleInteractor struct {
	RoleRepository RoleRepository
	Logger         api.Logger
}

// RoleRepository Persistence Interface
type RoleRepository interface {
	FindAll() []Role
	FindById(id uuid.UUID) (Role, error)
	FindByName(string) (Role, error)
	CreateOne(role Role) (Role, error)
	DeleteById(id uuid.UUID) error
}

// Role Entity
type Role struct {
	Id           uuid.UUID  `json:"id"`
	Name         string     `json:"name" binding:"required"`
	Permissions  []string   `json:"permissions" binding:"required"`
	DateCreated  *time.Time `json:"dateCreated,omitempty"`
	DateModified *time.Time `json:"dateModified,omitempty"`
}

func (r *Role) toDTO() {
	r.DateCreated = nil
	r.DateModified = nil
}

func (r *RoleInteractor) FindAll() []Role {
	return r.RoleRepository.FindAll()
}

func (r *RoleInteractor) FindById(id uuid.UUID) (Role, error) {
	return r.RoleRepository.FindById(id)
}

func (r *RoleInteractor) FindByName(name string) (Role, error) {
	return r.RoleRepository.FindByName(name)
}

func (r *RoleInteractor) CreateOne(role Role) (Role, error) {
	return r.RoleRepository.CreateOne(role)
}

func (r *RoleInteractor) DeleteById(id uuid.UUID) error {
	return r.RoleRepository.DeleteById(id)
}
