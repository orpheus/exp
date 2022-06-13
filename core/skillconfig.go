package core

// SkillConfigRepository defines the repository methods to
// interact with our entity.
type SkillConfigRepository interface {
	FindAll() []SkillConfig
	FindById(id string) (SkillConfig, error)
	CreateOne(skillConfig SkillConfig) (SkillConfig, error)
	CreateMany(skillConfigs []SkillConfig) error
	DeleteById(id string) error
}

// SkillConfig domain entity.
type SkillConfig struct {
	Id          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
