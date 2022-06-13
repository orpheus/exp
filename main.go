package main

import (
	"fmt"
	"github.com/orpheus/exp/infrastructure/postgres"
	"github.com/orpheus/exp/infrastructure/server"
	"os"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {
	dbUser := getEnv("DB_USER", "roark")
	dbPass := getEnv("DB_PASS", "")
	dbName := getEnv("DB_NAME", "exp")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")

	jdbcUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	conn := postgres.NewPgxPool(jdbcUrl)
	defer conn.Close()
	postgres.Migrate(conn)

	s := server.NewGin()
	server.Construct(s, conn)

	err := s.Run()
	if err != nil {
		panic("Error starting server")
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
