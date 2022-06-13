package repository

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/interfaces"
	usecases "github.com/orpheus/exp/usecases/auth"
	"log"
)

type RoleRepository struct {
	DB     PgxConn
	Logger interfaces.Logger
}

func (r *RoleRepository) FindAll() []usecases.Role {
	ds := "select * from role"
	rows, err := r.DB.Query(context.Background(), ds)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var records []usecases.Role
	for rows.Next() {
		var r usecases.Role
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

func (r *RoleRepository) FindById(id uuid.UUID) (usecases.Role, error) {
	sql := "select * from role where id = $1"
	var role usecases.Role
	err := r.DB.QueryRow(context.Background(), sql, id).
		Scan(&role.Id, &role.Name, &role.Permissions, &role.DateCreated, &role.DateModified)
	return role, err
}

func (r *RoleRepository) CreateOne(role usecases.Role) (usecases.Role, error) {
	ds := "insert into role (name, permissions) " +
		"VALUES ($1, $2) " +
		"RETURNING id, name, permissions, date_created, date_modified"

	var newRole usecases.Role
	err := r.DB.QueryRow(context.Background(), ds, role.Name, role.Permissions).
		Scan(&newRole.Id, &newRole.Name, &newRole.Permissions, &newRole.DateCreated, &newRole.DateModified)

	return newRole, err
}

func (r *RoleRepository) DeleteById(id uuid.UUID) error {
	ds := "delete from role where id = $1"
	_, err := r.DB.Exec(context.Background(), ds, id)
	return err
}
