package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/orpheus/exp/interfaces/ginhttprouter"
	"os"
	"time"
)

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {
	//dbUser := getEnv("DB_USER", "roark")
	//dbPass := getEnv("DB_PASS", "")
	//dbName := getEnv("DB_NAME", "exp")
	//dbHost := getEnv("DB_HOST", "localhost")
	//dbPort := getEnv("DB_PORT", "5432")
	//dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	//
	//dbConfig, err := pgxpool.ParseConfig(dbUrl)
	//if err != nil {
	//	log.Fatalln("Could not parse db connection string")
	//}
	//
	//conn, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	//	os.Exit(1)
	//}
	//defer conn.Close()
	//
	//dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	//	conn.ConnInfo().RegisterDataType(pgtype.DataType{
	//		Value: &pgtypeuuid.UUID{},
	//		Name:  "uuid",
	//		OID:   pgtype.UUIDOID,
	//	})
	//	return nil
	//}
	//
	//db.Migrate(conn)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST, OPTIONS, GET, PUT"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	ginhttprouter.RegisterRoutes(r)

	err := r.Run()
	if err != nil {
		panic("Error starting server")
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
