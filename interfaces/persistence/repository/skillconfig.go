package repository

//
//import (
//	"context"
//	"fmt"
//	"github.com/jackc/pgx/v4/pgxpool"
//	"log"
//	"strings"
//)
//
//
//
//type SkillConfigRepo struct {
//	DB *pgxpool.Pool
//}
//
//func (s *SkillConfigRepo) FindAll() []SkillConfig {
//	skillConfigRows, err := s.DB.Query(context.Background(), "select * from skill_config")
//	if err != nil {
//		panic(err)
//	}
//	defer skillConfigRows.Close()
//
//	var skillConfigs []SkillConfig
//	for skillConfigRows.Next() {
//		r := new(SkillConfig)
//		err := skillConfigRows.Scan(&r.Id, &r.Name, &r.Description)
//		if err != nil {
//			log.Fatal(err)
//		}
//		skillConfigs = append(skillConfigs, *r)
//	}
//	if err := skillConfigRows.Err(); err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(skillConfigs)
//	return skillConfigs
//}
//
//func (s *SkillConfigRepo) FindById(id string) (SkillConfig, error) {
//	r := new(SkillConfig)
//	err := s.DB.QueryRow(
//		context.Background(),
//		"select * from skill_config where id = $1", id).
//		Scan(&r.Id, &r.Name, &r.Description)
//	return *r, err
//}
//
//func (s *SkillConfigRepo) CreateOne(skillConfig SkillConfig) (SkillConfig, error) {
//	fmtStr := "insert into skill_config (id, name, description) VALUES ('%s', '%s', '%s') RETURNING id, name, description"
//	sql := fmt.Sprintf(fmtStr, skillConfig.Id, skillConfig.Name, skillConfig.Description)
//	r := new(SkillConfig)
//	err := s.DB.QueryRow(context.Background(), sql).
//		Scan(&r.Id, &r.Name, &r.Description)
//
//	fmt.Println(r, err)
//	return *r, err
//}
//
//func (s *SkillConfigRepo) CreateMany(skillConfigs []SkillConfig) (string, error) {
//	// Create query string for insert many
//	sqlBase := "insert into skill_config (id, name, description) VALUES %s RETURNING id, name, description"
//	var sqlValues []string
//	for _, v := range skillConfigs {
//		col := fmt.Sprintf("('%s', '%s', '%s')", v.Id, v.Name, v.Description)
//		sqlValues = append(sqlValues, col)
//	}
//	sqlValuesJoined := strings.Join(sqlValues, ", ")
//	sql := fmt.Sprintf(sqlBase, sqlValuesJoined)
//
//	// Query db
//	sqlResult, err := s.DB.Exec(context.Background(), sql)
//	if err != nil {
//		return sqlResult.String(), err
//	}
//	return sqlResult.String(), nil
//}
//
//func (s *SkillConfigRepo) DeleteById(id string) (string, error) {
//	sql := fmt.Sprintf("delete from skill_config where id = '%s'", id)
//	fmt.Println(sql)
//	_, err := s.DB.Query(context.Background(), sql)
//	return fmt.Sprintf("Deleted SkillConfig with id: %s", id), err
//}
