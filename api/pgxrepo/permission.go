package pgxrepo

import (
	"context"
	"fmt"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/system/sysauth"
	"log"
)

type PermissionRepository struct {
	DB     PgxConn
	Logger api.Logger
}

func (p *PermissionRepository) FindAll() ([]sysauth.Permission, error) {
	ds := "select * from permission"
	rows, err := p.DB.Query(context.Background(), ds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []sysauth.Permission
	for rows.Next() {
		var r sysauth.Permission
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

func (p *PermissionRepository) FindById(id string) (sysauth.Permission, error) {
	sql := "select * from permission where id = $1"
	var r sysauth.Permission
	err := p.DB.QueryRow(context.Background(), sql, id).
		Scan(&r.Id, &r.DateCreated)
	return r, err
}

func (p *PermissionRepository) CreateOne(id string) (sysauth.Permission, error) {
	ds := "insert into permission (id) VALUES ($1) RETURNING id, date_created"

	var r sysauth.Permission
	err := p.DB.QueryRow(context.Background(), ds, id).
		Scan(&r.Id, &r.DateCreated)

	return r, err
}

func (p *PermissionRepository) DeleteById(id string) error {
	ds := "delete from permission where id = $1"
	_, err := p.DB.Exec(context.Background(), ds, id)
	return err
}
