package pgxrepo

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/system/sysauth"
	"log"
)

type RoleRepository struct {
	DB     PgxConn
	Logger api.Logger
}

func (r *RoleRepository) FindAll() []sysauth.Role {
	ds := "select * from role"
	rows, err := r.DB.Query(context.Background(), ds)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var records []sysauth.Role
	for rows.Next() {
		var r sysauth.Role
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

func (r *RoleRepository) FindById(id uuid.UUID) (sysauth.Role, error) {
	sql := "select * from role where id = $1"
	var role sysauth.Role
	err := r.DB.QueryRow(context.Background(), sql, id).
		Scan(&role.Id, &role.Name, &role.Permissions, &role.DateCreated, &role.DateModified)
	return role, err
}

func (r *RoleRepository) FindByName(name string) (sysauth.Role, error) {
	sql := "select * from role where name = $1"
	var role sysauth.Role
	err := r.DB.QueryRow(context.Background(), sql, name).
		Scan(&role.Id, &role.Name, &role.Permissions, &role.DateCreated, &role.DateModified)
	return role, err
}

func (r *RoleRepository) CreateOne(role sysauth.Role) (sysauth.Role, error) {
	ds := "insert into role (name, permissions) " +
		"VALUES ($1, $2) " +
		"RETURNING id, name, permissions, date_created, date_modified"

	var newRole sysauth.Role
	err := r.DB.QueryRow(context.Background(), ds, role.Name, role.Permissions).
		Scan(&newRole.Id, &newRole.Name, &newRole.Permissions, &newRole.DateCreated, &newRole.DateModified)

	return newRole, err
}

func (r *RoleRepository) DeleteById(id uuid.UUID) error {
	ds := "delete from role where id = $1"
	_, err := r.DB.Exec(context.Background(), ds, id)
	return err
}
