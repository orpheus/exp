package main

import (
	"com.orpheus/exp/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
)

type Env struct {
	db *pgxpool.Pool
}

func main() {
	conn, err := pgxpool.Connect(context.Background(), "postgresql://roark:@localhost:5432/exp")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	r := gin.Default()
	env := &Env{db: conn}

	r.GET("/ping", ping)
	//r.GET("/skill", env.getSkill)
	r.GET("/skillConfig", env.getSkillConfig)
	//r.POST("/skillConfig", createSkillConfig)

	err = r.Run()
	if err != nil {
		panic("Error starting server")
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (e *Env) getSkill(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "OK")
}

func (e *Env) getSkillConfig(c *gin.Context) {
	skills := findAllSkills(e.db)
	c.IndentedJSON(http.StatusOK, skills)
}

func getSkillConfigById(c *gin.Context) {
	id := c.Param("id")

	c.IndentedJSON(http.StatusOK, id)
	//c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func createSkillConfig(c *gin.Context) {
	var skillConfig models.SkillConfig

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&skillConfig); err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, skillConfig)
}

func createSkill(c *gin.Context) {
	var skill models.Skill

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&skill); err != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, skill)
}

func findAllSkills(conn *pgxpool.Pool) []models.SkillConfig {
	skillConfigRows, err := conn.Query(context.Background(), "select * from skill_config")
	if err != nil {
		panic(err)
	}
	defer skillConfigRows.Close()

	var skillConfigs []models.SkillConfig
	for skillConfigRows.Next() {
		var r models.SkillConfig
		err := skillConfigRows.Scan(&r.Id, &r.Name, &r.Description)
		if err != nil {
			log.Fatal(err)
		}
		skillConfigs = append(skillConfigs, r)
	}
	if err := skillConfigRows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(skillConfigs)
	return skillConfigs
}
