package repository

import (
	"context"
	"fmt"
	usecases "github.com/orpheus/exp/usecases/auth"
	"log"
)

type PermissionRepo struct {
	DB PgxConn
}

func (svc *PermissionRepo) FindAll() ([]usecases.Permission, error) {
	ds := "select * from permission"
	rows, err := svc.DB.Query(context.Background(), ds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []usecases.Permission
	for rows.Next() {
		var r usecases.Permission
		err := rows.Scan(&r.Id, &r.DateCreated)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		records = append(records, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Printf("Fetched %d permissions\n", len(records))
	return records, nil
}

func (svc *PermissionRepo) FindById(id string) (usecases.Permission, error) {
	sql := "select * from permission where id = $1"
	var r usecases.Permission
	err := svc.DB.QueryRow(context.Background(), sql, id).
		Scan(&r.Id, &r.DateCreated)
	return r, err
}

func (svc *PermissionRepo) CreateOne(id string) (usecases.Permission, error) {
	ds := "insert into permission (id) VALUES ($1) RETURNING id, date_created"

	var r usecases.Permission
	err := svc.DB.QueryRow(context.Background(), ds, id).
		Scan(&r.Id, &r.DateCreated)

	return r, err
}

func (svc *PermissionRepo) DeleteById(id string) error {
	ds := "delete from permission where id = $1"
	_, err := svc.DB.Exec(context.Background(), ds, id)
	return err
}
