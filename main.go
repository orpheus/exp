package main

import (
	"com.orpheus/exp/controllers"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

func main() {
	conn, err := pgxpool.Connect(context.Background(), "postgresql://roark:@localhost:5432/exp")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	r := gin.Default()

	controllers.RegisterAll(r, conn)

	err = r.Run()
	if err != nil {
		panic("Error starting server")
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
