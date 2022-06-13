package repository

import (
	"context"
	"fmt"
	"github.com/orpheus/exp/domain"
	"log"
	"strings"
)

type SkillConfigRepo struct {
	DB PgxConn
}

func (s *SkillConfigRepo) FindAll() []domain.SkillConfig {
	skillConfigRows, err := s.DB.Query(context.Background(), "select * from skill_config")
	if err != nil {
		panic(err)
	}
	defer skillConfigRows.Close()

	var skillConfigs []domain.SkillConfig
	for skillConfigRows.Next() {
		r := new(domain.SkillConfig)
		err := skillConfigRows.Scan(&r.Id, &r.Name, &r.Description)
		if err != nil {
			log.Fatal(err)
		}
		skillConfigs = append(skillConfigs, *r)
	}
	if err := skillConfigRows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(skillConfigs)
	return skillConfigs
}

func (s *SkillConfigRepo) FindById(id string) (domain.SkillConfig, error) {
	r := new(domain.SkillConfig)
	err := s.DB.QueryRow(
		context.Background(),
		"select * from skill_config where id = $1", id).
		Scan(&r.Id, &r.Name, &r.Description)
	return *r, err
}

func (s *SkillConfigRepo) CreateOne(skillConfig domain.SkillConfig) (domain.SkillConfig, error) {
	fmtStr := "insert into skill_config (id, name, description) VALUES ('%s', '%s', '%s') RETURNING id, name, description"
	sql := fmt.Sprintf(fmtStr, skillConfig.Id, skillConfig.Name, skillConfig.Description)
	r := new(domain.SkillConfig)
	err := s.DB.QueryRow(context.Background(), sql).
		Scan(&r.Id, &r.Name, &r.Description)

	fmt.Println(r, err)
	return *r, err
}

func (s *SkillConfigRepo) CreateMany(skillConfigs []domain.SkillConfig) error {
	// Create query string for insert many
	sqlBase := "insert into skill_config (id, name, description) VALUES %s RETURNING id, name, description"
	var sqlValues []string
	for _, v := range skillConfigs {
		col := fmt.Sprintf("('%s', '%s', '%s')", v.Id, v.Name, v.Description)
		sqlValues = append(sqlValues, col)
	}
	sqlValuesJoined := strings.Join(sqlValues, ", ")
	sql := fmt.Sprintf(sqlBase, sqlValuesJoined)

	// Query db
	_, err := s.DB.Exec(context.Background(), sql)
	if err != nil {
		return err
	}
	return nil
}

func (s *SkillConfigRepo) DeleteById(id string) error {
	sql := fmt.Sprintf("delete from skill_config where id = '%s'", id)
	fmt.Println(sql)
	_, err := s.DB.Query(context.Background(), sql)
	return err
}
