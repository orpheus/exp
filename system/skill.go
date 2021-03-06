package system

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/core"
	"log"
	"time"
)

// SkillInteractor Service
type SkillInteractor struct {
	SkillRepository core.SkillRepository
	Policy          core.SkillPolicy
	Logger          api.Logger
}

func (s *SkillInteractor) FindAllSkills() []core.Skill {
	skills, err := s.SkillRepository.FindAll()
	if err != nil {
		message := fmt.Errorf("%s", err.Error())
		s.Logger.Log(message.Error())
		return nil
	}
	return skills
}

func (s *SkillInteractor) FindSkillById(id uuid.UUID) (core.Skill, error) {
	skill, err := s.SkillRepository.FindById(id)
	if err != nil {
		err = fmt.Errorf("%s", err.Error())
		s.Logger.Log(err.Error())
		return skill, err
	}
	return skill, nil
}

func (s *SkillInteractor) CreateSkill(skill core.Skill, userId uuid.UUID) (core.Skill, error) {
	exists, err := s.SkillRepository.ExistsBySkillIdAndUserId(skill.SkillId, userId)
	if err != nil {
		message := fmt.Errorf("%s", err.Error())
		s.Logger.Log(message.Error())
	}

	if exists {
		return skill, errors.New("skill already exists")
	}

	skill.UserId = userId

	savedSkill, err := s.SkillRepository.CreateOne(skill)
	if err != nil {
		return skill, err
	}
	return savedSkill, nil
}

func (s *SkillInteractor) DeleteById(skillId uuid.UUID) error {
	_, err := s.SkillRepository.FindById(skillId)
	if err != nil {
		return err
	}
	_, err = s.SkillRepository.DeleteById(skillId)
	if err != nil {
		return err
	}
	return nil
}

func (s *SkillInteractor) ExistsBySkillIdAndUserId(skillId string, userId uuid.UUID) (bool, error) {
	return s.SkillRepository.ExistsBySkillIdAndUserId(skillId, userId)
}

func (s *SkillInteractor) AddTxp(txp int, skillId uuid.UUID) (*core.Skill, error) {
	skill, err := s.FindSkillById(skillId)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	if txp <= 0 {
		return nil, errors.New("cannot add txp less than or equal to 0")
	}

	now := time.Now()
	last := skill.DateLastTxpAdd
	secondsSinceLastUpdate := int(now.Sub(last).Seconds())
	allowFirstTimeTxpAdd := s.Policy.AllowFirstTimeTxpAdd(skill, txp)

	if txp > secondsSinceLastUpdate && !allowFirstTimeTxpAdd {
		var message string
		if skill.IsNewSkill() {
			message = "Cannot add more than 3600 txp for the first hour of the skill's lifetime"
		} else {
			message = "Cannot add more txp than the difference of time in seconds between now and the last update"
		}
		return nil, errors.New(message)
	}

	skill.AddTxp(txp)

	skill, err = s.SkillRepository.UpdateExpLvl(skill)
	return &skill, err
}
