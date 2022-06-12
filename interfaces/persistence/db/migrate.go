package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

func Migrate(conn *pgxpool.Pool) {
	files, err := ioutil.ReadDir("db/sql")
	if err != nil {
		log.Fatalf("Failed to read db/sql dir in migarte: %s", err.Error())
	}

	tx, err := conn.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		log.Fatalln("Failed to create TX in migration")
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
			fmt.Println("Rolled back migration")
		} else {
			tx.Commit(context.TODO())
			fmt.Println("Committed migrations")
		}
	}()

	versions := make(map[string]string)
	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
		fn := file.Name()
		validator := regexp.MustCompile("^V\\d+.\\d+.\\d+__\\w+\\.(sql)$")
		validMigrationScript := validator.MatchString(fn)
		if !validMigrationScript {
			log.Fatalf("Invalid migration script format: %s. Expecting: V<X>.<X>.<X>__<NAME>.sql", fn)
		}

		split := strings.Split(fn, "__")
		version := split[0]
		if _, ok := versions[version]; ok {
			log.Fatalln("Duplicate versions found in migration sql scripts")
		}
		versions[version] = fn

		c, ioErr := ioutil.ReadFile(fmt.Sprintf("db/sql/%s", fn))
		if ioErr != nil {
			log.Fatalf("Could not read sql file: %s", fn)
		}
		sql := string(c)
		_, err = conn.Exec(context.Background(), sql)
		if err != nil {
			log.Fatalf("Failed to execute migartion script: %s", err.Error())
		}
	}
}
