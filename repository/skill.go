package repository

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/orpheus/exp/core"
	"log"
	"time"
)

type SkillRepo struct {
	DB *pgxpool.Pool
}

func (s *SkillRepo) FindAll() []core.Skill {
	SkillRows, err := s.DB.Query(context.Background(), "select * from skill")
	if err != nil {
		panic(err)
	}
	defer SkillRows.Close()

	var skills []core.Skill
	for SkillRows.Next() {
		var r core.Skill
		err := SkillRows.Scan(&r.Id, &r.SkillId, &r.UserId, &r.Exp, &r.Txp, &r.Level, &r.DateCreated, &r.DateModified)
		if err != nil {
			log.Fatal(err)
		}
		skills = append(skills, r)
	}
	if err := SkillRows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetched %d skills\n", len(skills))
	return skills
}

func (s *SkillRepo) FindById(id uuid.UUID) (core.Skill, error) {
	sql := "select * from skill where id = $1"
	var skill core.Skill
	err := s.DB.QueryRow(context.Background(), sql, id).
		Scan(&skill.Id, &skill.SkillId, &skill.UserId, &skill.Exp, &skill.Txp, &skill.Level, &skill.DateCreated, &skill.DateModified)
	fmt.Println(skill)
	return skill, err
}

func (s *SkillRepo) CreateOne(skill core.Skill) (core.Skill, error) {
	sql := "insert into skill (skill_id, user_id, exp, txp, level, date_created, date_modified) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7) " +
		"RETURNING id, skill_id, user_id, exp, txp, level, date_created, date_modified"

	var sk core.Skill
	// TODO(Make level 1 by default in postgres)
	err := s.DB.QueryRow(context.Background(), sql, skill.SkillId, nil, skill.Exp, skill.Txp, 1, time.Now(), time.Now()).
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

func (s *SkillRepo) UpdateExpLvl(skill core.Skill) (core.Skill, error) {
	sql := "update skill set exp = $1, txp = $2, level = $3, date_modified = $4 where id = $5 " +
		"RETURNING id, skill_id, user_id, exp, txp, level, date_created, date_modified"

	var sk core.Skill
	err := s.DB.QueryRow(context.Background(), sql, skill.Exp, skill.Txp, skill.Level, time.Now(), skill.Id).
		Scan(&sk.Id, &sk.SkillId, &sk.UserId, &sk.Exp, &sk.Txp, &sk.Level, &sk.DateCreated, &sk.DateModified)
	return sk, err
}
