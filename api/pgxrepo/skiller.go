package pgxrepo

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/core"
	"log"
	"time"
)

type SkillerRepository struct {
	DB     PgxConn
	Logger api.Logger
}

func (r *SkillerRepository) FindAll() ([]core.Skiller, error) {
	userRows, err := r.DB.Query(context.Background(), "select * from user_account")
	if err != nil {
		return nil, err
	}
	defer userRows.Close()

	var users []core.Skiller
	for userRows.Next() {
		u := new(core.Skiller)
		err := userRows.Scan(&u.Id, &u.Username, &u.Password, &u.Email, &u.RoleId, &u.DateCreated, &u.DateModified)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, *u)
	}
	if err := userRows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *SkillerRepository) FindByUsername(username string) (core.Skiller, error) {
	ds := "select * from user_account where username = $1"

	var u core.Skiller
	err := r.DB.QueryRow(context.Background(), ds, username).
		Scan(&u.Id, &u.Username, &u.Password, &u.Email, &u.RoleId, &u.DateCreated, &u.DateModified)

	if err != nil {
		err = fmt.Errorf("db scan failed for Skiller in repository.skiller.FindByUsername: %s", err.Error())
		r.Logger.Log(err.Error())
		return core.Skiller{}, err
	}

	return u, nil
}

func (r *SkillerRepository) FindByEmail(email string) (core.Skiller, error) {
	ds := "select * from user_account where email = $1"

	var u core.Skiller
	err := r.DB.QueryRow(context.Background(), ds, email).
		Scan(&u.Id, &u.Username, &u.Email, &u.RoleId, &u.DateCreated, &u.DateModified)

	r.Logger.Log(fmt.Sprintf("%s", u))

	return u, err
}

func (r *SkillerRepository) FindById(id uuid.UUID) (core.Skiller, error) {
	ds := "select * from user_account where id = $1"

	var u core.Skiller
	err := r.DB.QueryRow(context.Background(), ds, id).
		Scan(&u.Id, &u.Username, &u.Email, &u.RoleId, &u.DateCreated, &u.DateModified)

	r.Logger.Log(fmt.Sprintf("%s", u))

	return u, err
}

// CreateOne TODO(Handle Email)
func (r *SkillerRepository) CreateOne(skiller core.Skiller) (core.Skiller, error) {
	ds := "insert into user_account (username, password, email, role_id, date_created, date_modified) " +
		"VALUES($1, $2, $3, $4, $5, $6) " +
		"RETURNING id, username, email, role_id, date_created, date_modified"

	var u core.Skiller
	err := r.DB.QueryRow(context.Background(), ds, skiller.Username, skiller.Password, skiller.Email, skiller.RoleId, time.Now(), time.Now()).
		Scan(&u.Id, &u.Username, &u.Email, &u.RoleId, &u.DateCreated, &u.DateModified)

	return u, err
}
