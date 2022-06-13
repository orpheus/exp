package pgxrepo

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/core"
	"log"
	"time"
)

type SkillRepository struct {
	DB     PgxConn
	Logger api.Logger
}

func (s *SkillRepository) FindAll() ([]core.Skill, error) {
	rows, err := s.DB.Query(context.Background(), "select * from skill")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []core.Skill
	for rows.Next() {
		var r core.Skill
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
		return []core.Skill{}, nil
	}
	return skills, nil
}

func (s *SkillRepository) FindById(id uuid.UUID) (core.Skill, error) {
	sql := "select * from skill where id = $1"
	var skill core.Skill
	err := s.DB.QueryRow(context.Background(), sql, id).
		Scan(&skill.Id, &skill.SkillId, &skill.UserId, &skill.Exp, &skill.Txp, &skill.Level, &skill.DateCreated, &skill.DateModified, &skill.DateLastTxpAdd)
	return skill, err
}

func (s *SkillRepository) CreateOne(skill core.Skill) (core.Skill, error) {
	sql := "insert into skill (skill_id, user_id, exp, txp, level, date_created, date_modified) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7) " +
		"RETURNING id, skill_id, user_id, exp, txp, level, date_created, date_modified, date_last_txp_add"

	var sk core.Skill
	// TODO(Make level 1 by default in postgres)
	err := s.DB.QueryRow(context.Background(), sql, skill.SkillId, skill.UserId, skill.Exp, skill.Txp, 1, time.Now(), time.Now()).
		Scan(&sk.Id, &sk.SkillId, &sk.UserId, &sk.Exp, &sk.Txp, &sk.Level, &sk.DateCreated, &sk.DateModified, &sk.DateLastTxpAdd)

	return sk, err
}

func (s *SkillRepository) DeleteById(id uuid.UUID) (string, error) {
	// TODO(Check if exists first, so you can let client know he did what was expected)
	sql := "delete from skill where id = $1"
	_, err := s.DB.Exec(context.Background(), sql, id)
	return fmt.Sprintf("Deleted Skill with id: %v", id), err
}

func (s *SkillRepository) UpdateExpLvl(skill core.Skill) (core.Skill, error) {
	sql := "update skill set exp = $1, txp = $2, level = $3, date_modified = $4, date_last_txp_add = $5 where id = $6 " +
		"RETURNING id, skill_id, user_id, exp, txp, level, date_created, date_modified, date_last_txp_add"

	var sk core.Skill
	err := s.DB.QueryRow(context.Background(), sql, skill.Exp, skill.Txp, skill.Level, time.Now(), time.Now(), skill.Id).
		Scan(&sk.Id, &sk.SkillId, &sk.UserId, &sk.Exp, &sk.Txp, &sk.Level, &sk.DateCreated, &sk.DateModified, &sk.DateLastTxpAdd)
	return sk, err
}

func (s *SkillRepository) ExistsBySkillIdAndUserId(skillId string, userId uuid.UUID) (bool, error) {
	ds := "select exists (select true from skill where skill_id=$1 and user_id=$2)"
	var b bool
	err := s.DB.QueryRow(context.Background(), ds, skillId, userId).
		Scan(&b)
	return b, err
}
