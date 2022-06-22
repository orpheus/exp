package system

import (
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/core"
)

type SkillerInteractor struct {
	Repo   core.SkillerRepository
	Logger api.Logger
}

func (s *SkillerInteractor) FindAll() ([]core.Skiller, error) {
	return s.Repo.FindAll()
}

func (s *SkillerInteractor) FindById(id uuid.UUID) (core.Skiller, error) {
	return s.Repo.FindById(id)
}

func (s *SkillerInteractor) FindByEmail(email string) (core.Skiller, error) {
	return s.Repo.FindByEmail(email)
}

func (s *SkillerInteractor) FindByUsername(username string) (core.Skiller, error) {
	return s.Repo.FindByUsername(username)
}

func (s *SkillerInteractor) CreateOne(skiller core.Skiller) (core.Skiller, error) {
	return s.Repo.CreateOne(skiller)
}
