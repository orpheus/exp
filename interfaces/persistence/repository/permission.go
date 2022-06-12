package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type Permission struct {
	Id          string    `json:"id"`
	DateCreated time.Time `json:"dateCreated"`
}

type PermissionRepo struct {
	DB *pgxpool.Pool
}

func (svc *PermissionRepo) FindAll() ([]Permission, error) {
	ds := "select * from permission"
	rows, err := svc.DB.Query(context.Background(), ds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []Permission
	for rows.Next() {
		var r Permission
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

func (svc *PermissionRepo) FindById(id string) (Permission, error) {
	sql := "select * from permission where id = $1"
	var r Permission
	err := svc.DB.QueryRow(context.Background(), sql, id).
		Scan(&r.Id, &r.DateCreated)
	return r, err
}

func (svc *PermissionRepo) CreateOne(id string) (Permission, error) {
	ds := "insert into permission (id) VALUES ($1) RETURNING id, date_created"

	var r Permission
	err := svc.DB.QueryRow(context.Background(), ds, id).
		Scan(&r.Id, &r.DateCreated)

	return r, err
}

func (svc *PermissionRepo) DeleteById(id string) (string, error) {
	ds := "delete from permission where id = $1"
	_, err := svc.DB.Exec(context.Background(), ds, id)
	return fmt.Sprintf("Deleted Permission with id: %v", id), err
}
