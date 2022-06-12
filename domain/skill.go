package domain

import (
	"github.com/gofrs/uuid"
	"time"
)

type SkillRepository interface {
	FindAll() ([]Skill, error)
	FindById(id uuid.UUID) (Skill, error)
	CreateOne(skill Skill) (Skill, error)
	DeleteById(id uuid.UUID) (string, error)
	UpdateExpLvl(skill Skill) (Skill, error)
	ExistsByUserId(skillId uuid.UUID, userId uuid.UUID) (bool, error)
}

type Skill struct {
	Id             uuid.UUID `json:"id"`
	SkillId        string    `json:"skillId" binding:"required"`
	UserId         uuid.UUID `json:"userId"`
	Exp            int       `json:"exp"`
	Txp            int       `json:"txp"`
	Level          int       `json:"level"`
	DateCreated    time.Time `json:"dateCreated"`
	DateModified   time.Time `json:"dateModified"`
	DateLastTxpAdd time.Time `json:"dateLastTxpAdd"`
}

type SkillPolicy struct{}

var skillTable = createSkillTable()

func (s *Skill) AddTxp(txp int) {
	s.Txp += txp
	s.Exp += txp
	s.tryLevelIncrease()
}

func (s *Skill) tryLevelIncrease() {
	for i := 2; i <= 99; i++ {
		if s.Exp >= skillTable[i] {
			if s.Level < i {
				s.Level = i
			}
		} else {
			break
		}
	}
}

func createSkillTable() map[int]int {
	const rate = 1.0777
	exp := 3600
	diff := exp
	level := 2

	m := map[int]int{1: 0}

	for level <= 99 {
		m[level] = exp

		diff = int(float64(diff) * rate)
		exp = exp + diff
		level += 1
	}

	return m
}

// FirstTimeTxp checks to see if there is any txp on the skill and if
// not that the requested txp to add to the skill is less than one hour.
// Policy: Users can only add up to an hour of txp their first time.
func (p *SkillPolicy) FirstTimeTxp(skill Skill, requestedTxp int) bool {
	now := time.Now()
	if skill.Txp == 0 && now.Sub(skill.DateCreated).Hours() < 1 {
		if requestedTxp <= int(time.Hour.Seconds()) {
			return true
		}
	}
	return false
}
