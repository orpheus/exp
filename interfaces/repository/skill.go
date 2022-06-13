package repository

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/domain"
	"log"
	"time"
)

type SkillRepo struct {
	DB PgxConn
}

func (s *SkillRepo) FindAll() ([]domain.Skill, error) {
	rows, err := s.DB.Query(context.Background(), "select * from skill")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []domain.Skill
	for rows.Next() {
		var r domain.Skill
		err := rows.Scan(&r.Id, &r.SkillId, &r.UserId, &r.Exp, &r.Txp, &r.Level, &r.DateCreated, &r.DateModified, &r.DateLastTxpAdd)
		if err != nil {
			log.Fatal(err)
		}
		skills = append(skills, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetched %d skills\n", len(skills))
	if len(skills) == 0 {
		return []domain.Skill{}, nil
	}
	return skills, nil
}

func (s *SkillRepo) FindById(id uuid.UUID) (domain.Skill, error) {
	sql := "select * from skill where id = $1"
	var skill domain.Skill
	err := s.DB.QueryRow(context.Background(), sql, id).
		Scan(&skill.Id, &skill.SkillId, &skill.UserId, &skill.Exp, &skill.Txp, &skill.Level, &skill.DateCreated, &skill.DateModified, &skill.DateLastTxpAdd)
	return skill, err
}

func (s *SkillRepo) CreateOne(skill domain.Skill) (domain.Skill, error) {
	sql := "insert into skill (skill_id, user_id, exp, txp, level, date_created, date_modified) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7) " +
		"RETURNING id, skill_id, user_id, exp, txp, level, date_created, date_modified, date_last_txp_add"

	var sk domain.Skill
	// TODO(Make level 1 by default in postgres)
	err := s.DB.QueryRow(context.Background(), sql, skill.SkillId, skill.UserId, skill.Exp, skill.Txp, 1, time.Now(), time.Now()).
		Scan(&sk.Id, &sk.SkillId, &sk.UserId, &sk.Exp, &sk.Txp, &sk.Level, &sk.DateCreated, &sk.DateModified, &sk.DateLastTxpAdd)

	return sk, err
}

func (s *SkillRepo) DeleteById(id uuid.UUID) (string, error) {
	// TODO(Check if exists first, so you can let client know he did what was expected)
	sql := "delete from skill where id = $1"
	_, err := s.DB.Exec(context.Background(), sql, id)
	return fmt.Sprintf("Deleted Skill with id: %v", id), err
}

func (s *SkillRepo) UpdateExpLvl(skill domain.Skill) (domain.Skill, error) {
	sql := "update skill set exp = $1, txp = $2, level = $3, date_modified = $4, date_last_txp_add = $5 where id = $6 " +
		"RETURNING id, skill_id, user_id, exp, txp, level, date_created, date_modified, date_last_txp_add"

	var sk domain.Skill
	err := s.DB.QueryRow(context.Background(), sql, skill.Exp, skill.Txp, skill.Level, time.Now(), time.Now(), skill.Id).
		Scan(&sk.Id, &sk.SkillId, &sk.UserId, &sk.Exp, &sk.Txp, &sk.Level, &sk.DateCreated, &sk.DateModified, &sk.DateLastTxpAdd)
	return sk, err
}

func (s *SkillRepo) ExistsByUserId(skillId string, userId uuid.UUID) (bool, error) {
	ds := "select exists (select true from skill where skill_id=$1 and user_id=$2)"
	var b bool
	err := s.DB.QueryRow(context.Background(), ds, skillId, userId).
		Scan(&b)
	return b, err
}
