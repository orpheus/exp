package core

import (
	"github.com/gofrs/uuid"
	"time"
)

type Skill struct {
	Id           uuid.UUID `json:"id"`
	SkillId      string    `json:"skillId" binding:"required"`
	UserId       uuid.UUID `json:"userId"`
	Exp          int       `json:"exp"`
	Txp          int       `json:"txp"`
	Level        int       `json:"level"`
	DateCreated  time.Time `json:"dateCreated"`
	DateModified time.Time `json:"dateModified"`
}

var skillTable = createSkillTable()

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
