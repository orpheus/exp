package repository

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/domain"
	"github.com/orpheus/exp/interfaces"
	"log"
	"time"
)

type SkillerRepo struct {
	DB     PgxConn
	Logger interfaces.Logger
}

func (r *SkillerRepo) FindAll() ([]domain.Skiller, error) {
	userRows, err := r.DB.Query(context.Background(), "select * from user_account")
	if err != nil {
		return nil, err
	}
	defer userRows.Close()

	var users []domain.Skiller
	for userRows.Next() {
		u := new(domain.Skiller)
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

func (r *SkillerRepo) FindByUsername(username string) (domain.Skiller, error) {
	ds := "select * from user_account where username = $1"

	var u domain.Skiller
	err := r.DB.QueryRow(context.Background(), ds, username).
		Scan(&u.Id, &u.Username, &u.Email, &u.RoleId, &u.DateCreated, &u.DateModified)

	r.Logger.Log(fmt.Sprintf("%s", u))

	return u, err
}

func (r *SkillerRepo) FindByEmail(email string) (domain.Skiller, error) {
	ds := "select * from user_account where email = $1"

	var u domain.Skiller
	err := r.DB.QueryRow(context.Background(), ds, email).
		Scan(&u.Id, &u.Username, &u.Email, &u.RoleId, &u.DateCreated, &u.DateModified)

	r.Logger.Log(fmt.Sprintf("%s", u))

	return u, err
}

func (r *SkillerRepo) FindById(id uuid.UUID) (domain.Skiller, error) {
	ds := "select * from user_account where id = $1"

	var u domain.Skiller
	err := r.DB.QueryRow(context.Background(), ds, id).
		Scan(&u.Id, &u.Username, &u.Email, &u.RoleId, &u.DateCreated, &u.DateModified)

	r.Logger.Log(fmt.Sprintf("%s", u))

	return u, err
}

// CreateOne TODO(Handle Email)
func (r *SkillerRepo) CreateOne(user domain.Skiller, hashedPassword string) (domain.Skiller, error) {
	ds := "insert into user_account (username, password, email, role_id, date_created, date_modified) " +
		"VALUES($1, $2, $3, $4, $5, $6) " +
		"RETURNING id, username, email, role_id, date_created, date_modified"

	var u domain.Skiller
	err := r.DB.QueryRow(context.Background(), ds, user.Username, hashedPassword, user.Email, user.RoleId, time.Now(), time.Now()).
		Scan(&u.Id, &u.Username, &u.Email, &u.RoleId, &u.DateCreated, &u.DateModified)

	return u, err
}
