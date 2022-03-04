package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type SkillConfig struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SkillConfigRepo struct {
	DB *pgxpool.Pool
}

func (s *SkillConfigRepo) FindAll() []SkillConfig {
	skillConfigRows, err := s.DB.Query(context.Background(), "select * from skill_config")
	if err != nil {
		panic(err)
	}
	defer skillConfigRows.Close()

	var skillConfigs []SkillConfig
	for skillConfigRows.Next() {
		r := new(SkillConfig)
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
