package repository

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type Role struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name" binding:"required"`
	Permissions  []string  `json:"permissions" binding:"required"`
	DateCreated  time.Time `json:"dateCreated"`
	DateModified time.Time `json:"dateModified"`
}

type RoleRepo struct {
	DB *pgxpool.Pool
}

func (svc *RoleRepo) FindAll() []Role {
	ds := "select * from role"
	rows, err := svc.DB.Query(context.Background(), ds)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var records []Role
	for rows.Next() {
		var r Role
		err := rows.Scan(&r.Id, &r.Name, &r.Permissions, &r.DateCreated, &r.DateModified)
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, r)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetched %d roles\n", len(records))
	return records
}

func (svc *RoleRepo) FindById(id uuid.UUID) (Role, error) {
	sql := "select * from role where id = $1"
	var r Role
	err := svc.DB.QueryRow(context.Background(), sql, id).
		Scan(&r.Id, &r.Name, &r.Permissions, &r.DateCreated, &r.DateModified)
	return r, err
}

func (svc *RoleRepo) CreateOne(role Role) (Role, error) {
	ds := "insert into role (name, permissions) " +
		"VALUES ($1, $2) " +
		"RETURNING id, name, permissions, date_created, date_modified"

	var r Role
	err := svc.DB.QueryRow(context.Background(), ds, role.Name, role.Permissions).
		Scan(&r.Id, &r.Name, &r.Permissions, &r.DateCreated, &r.DateModified)

	return r, err
}

func (svc *RoleRepo) DeleteById(id uuid.UUID) (string, error) {
	ds := "delete from role where id = $1"
	_, err := svc.DB.Exec(context.Background(), ds, id)
	return fmt.Sprintf("Deleted Role with id: %v", id), err
}
