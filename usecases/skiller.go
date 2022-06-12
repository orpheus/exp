package usecases

import (
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/domain"
	"github.com/orpheus/exp/interfaces"
)

type SkillerInteractor struct {
	Repo   domain.SkillerRepository
	Logger interfaces.Logger
}

func (s *SkillerInteractor) FindAll() ([]domain.Skiller, error) {
	return s.Repo.FindAll()
}

func (s *SkillerInteractor) FindById(id uuid.UUID) (domain.Skiller, error) {
	return s.Repo.FindById(id)
}

func (s *SkillerInteractor) FindByEmail(email string) (domain.Skiller, error) {
	return s.Repo.FindByEmail(email)
}

func (s *SkillerInteractor) FindByUsername(username string) (domain.Skiller, error) {
	return s.Repo.FindByUsername(username)
}

func (s *SkillerInteractor) CreateOne(skiller domain.Skiller) (domain.Skiller, error) {
	// TODO(check if already exists)
	return s.Repo.CreateOne(skiller)
}
