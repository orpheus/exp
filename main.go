package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/orpheus/exp/controllers"
	"os"
)

func main() {
	dbConnString := "postgresql://roark:@localhost:5432/exp"
	dbConfig, err := pgxpool.ParseConfig(dbConnString)
	if err != nil {
		panic("Could not parse db connection string")
	}
	conn, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &pgtypeuuid.UUID{},
			Name:  "uuid",
			OID:   pgtype.UUIDOID,
		})
		return nil
	}
	defer conn.Close()

	r := gin.Default()
	controllers.RegisterAll(r, conn)

	err = r.Run()
	if err != nil {
		panic("Error starting server")
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
