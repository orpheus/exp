package system

import (
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/core"
)

// SkillConfigInteractor Service.
type SkillConfigInteractor struct {
	Repo   core.SkillConfigRepository
	Logger api.Logger
}

func (s *SkillConfigInteractor) FindAllSkillConfigs() []core.SkillConfig {
	return s.Repo.FindAll()
}

func (s *SkillConfigInteractor) FindSkillConfigById(id string) (core.SkillConfig, error) {
	return s.Repo.FindById(id)
}

func (s *SkillConfigInteractor) CreateSkillConfig(skillConfig core.SkillConfig) (core.SkillConfig, error) {
	return s.Repo.CreateOne(skillConfig)
}

func (s *SkillConfigInteractor) CreateSkillConfigs(skillConfigs []core.SkillConfig) error {
	return s.Repo.CreateMany(skillConfigs)
}

func (s *SkillConfigInteractor) DeleteById(id string) error {
	return s.Repo.DeleteById(id)
}
