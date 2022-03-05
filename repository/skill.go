package repository

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type Skill struct {
	Id           uuid.UUID     `json:"id"`
	SkillId      string        `json:"skillId" binding:"required"`
	UserId       uuid.NullUUID `json:"userId"`
	Exp          int           `json:"exp"`
	Txp          int           `json:"txp"`
	Level        int           `json:"level"`
	DateCreated  time.Time     `json:"dateCreated"`
	DateModified time.Time     `json:"dateModified"`
}

type SkillRepo struct {
	DB *pgxpool.Pool
}

func (s *SkillRepo) FindAll() []Skill {
	SkillRows, err := s.DB.Query(context.Background(), "select * from skill")
	if err != nil {
		panic(err)
	}
	defer SkillRows.Close()

	var skills []Skill
	for SkillRows.Next() {
		r := new(Skill)
		err := SkillRows.Scan(&r.Id, &r.SkillId, &r.UserId, &r.Exp, &r.Txp, &r.Level, &r.DateCreated, &r.DateModified)
		if err != nil {
			log.Fatal(err)
		}
		skills = append(skills, *r)
	}
	if err := SkillRows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetched %d skills\n", len(skills))
	return skills
}

func (s *SkillRepo) FindById(id uuid.UUID) (Skill, error) {
	sql := "select * from skill where id = $1"
	var skill Skill
	err := s.DB.QueryRow(context.Background(), sql, id).
		Scan(&skill.Id, &skill.SkillId, &skill.UserId, &skill.Exp, &skill.Txp, &skill.Level, &skill.DateCreated, &skill.DateModified)
	fmt.Println(skill)
	return skill, err
}

func (s *SkillRepo) CreateOne(skill Skill) (Skill, error) {
	sql := "insert into skill (skill_id, user_id, exp, txp, level, date_created, date_modified) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7) " +
		"RETURNING id, skill_id, user_id, exp, txp, level, date_created, date_modified"

	var sk Skill
	err := s.DB.QueryRow(context.Background(), sql, skill.SkillId, nil, skill.Exp, skill.Txp, skill.Level, time.Now(), time.Now()).
		Scan(&sk.Id, &sk.SkillId, &sk.UserId, &sk.Exp, &sk.Txp, &sk.Level, &sk.DateCreated, &sk.DateModified)

	fmt.Println(sk)

	return sk, err
}

func (s *SkillRepo) DeleteById(id uuid.UUID) (string, error) {
	// TODO(Check if exists first, so you can let client know he did what was expected)
	sql := "delete from skill where id = $1"
	_, err := s.DB.Exec(context.Background(), sql, id)
	return fmt.Sprintf("Deleted Skill with id: %v", id), err
}
