package usecases

import (
	"github.com/orpheus/exp/domain"
	"github.com/orpheus/exp/interfaces"
)

// SkillConfigInteractor Service.
type SkillConfigInteractor struct {
	Repo   domain.SkillConfigRepository
	Logger interfaces.Logger
}

func (s *SkillConfigInteractor) FindAllSkillConfigs() []domain.SkillConfig {
	return s.Repo.FindAll()
}

func (s *SkillConfigInteractor) FindSkillConfigById(id string) (domain.SkillConfig, error) {
	return s.Repo.FindById(id)
}

func (s *SkillConfigInteractor) CreateSkillConfig(skillConfig domain.SkillConfig) (domain.SkillConfig, error) {
	return s.Repo.CreateOne(skillConfig)
}

func (s *SkillConfigInteractor) CreateSkillConfigs(skillConfigs []domain.SkillConfig) error {
	return s.Repo.CreateMany(skillConfigs)
}

func (s *SkillConfigInteractor) DeleteById(id string) error {
	return s.Repo.DeleteById(id)
}
