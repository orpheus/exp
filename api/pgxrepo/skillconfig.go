package pgxrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/core"
	"log"
	"strings"
)

type SkillConfigRepository struct {
	DB     PgxConn
	Logger api.Logger
}

func (s *SkillConfigRepository) FindAll() ([]core.SkillConfig, error) {
	skillConfigRows, err := s.DB.Query(context.Background(), "select * from skill_config")
	if err != nil {
		return nil, err
	}
	defer skillConfigRows.Close()

	var skillConfigs []core.SkillConfig
	for skillConfigRows.Next() {
		r := new(core.SkillConfig)
		err := skillConfigRows.Scan(&r.Id, &r.Name, &r.Description)
		if err != nil {
			log.Fatal(err)
		}
		skillConfigs = append(skillConfigs, *r)
	}

	if err := skillConfigRows.Err(); err != nil {
		return nil, err
	}

	if len(skillConfigs) == 0 {
		return []core.SkillConfig{}, nil
	}

	return skillConfigs, nil
}

func (s *SkillConfigRepository) FindById(id string) (core.SkillConfig, error) {
	r := new(core.SkillConfig)
	err := s.DB.QueryRow(
		context.Background(),
		"select * from skill_config where id = $1", id).
		Scan(&r.Id, &r.Name, &r.Description)
	if err == pgx.ErrNoRows {
		return *r, errors.New(fmt.Sprintf("skill (%s) not found", id))
	}
	return *r, err
}

func (s *SkillConfigRepository) CreateOne(skillConfig core.SkillConfig) (core.SkillConfig, error) {
	fmtStr := "insert into skill_config (id, name, description) VALUES ('%s', '%s', '%s') RETURNING id, name, description"
	sql := fmt.Sprintf(fmtStr, skillConfig.Id, skillConfig.Name, skillConfig.Description)
	r := new(core.SkillConfig)
	err := s.DB.QueryRow(context.Background(), sql).
		Scan(&r.Id, &r.Name, &r.Description)

	fmt.Println(r, err)
	return *r, err
}

func (s *SkillConfigRepository) CreateMany(skillConfigs []core.SkillConfig) error {
	// Create query string for insert many
	sqlBase := "insert into skill_config (id, name, description) VALUES %s RETURNING id, name, description"
	var sqlValues []string
	for _, v := range skillConfigs {
		col := fmt.Sprintf("('%s', '%s', '%s')", v.Id, v.Name, v.Description)
		sqlValues = append(sqlValues, col)
	}
	sqlValuesJoined := strings.Join(sqlValues, ", ")
	sql := fmt.Sprintf(sqlBase, sqlValuesJoined)

	// Query pgx
	_, err := s.DB.Exec(context.Background(), sql)
	if err != nil {
		return err
	}
	return nil
}

func (s *SkillConfigRepository) DeleteById(id string) error {
	sql := fmt.Sprintf("delete from skill_config where id = '%s'", id)
	fmt.Println(sql)
	_, err := s.DB.Query(context.Background(), sql)
	return err
}
