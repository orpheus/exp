package repository

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type User struct {
	Id           uuid.UUID     `json:"id"`
	Username     string        `json:"username" binding:"required"`
	Password     string        `json:"password,omitempty"`
	Email        string        `json:"email,omitempty"`
	RoleId       uuid.NullUUID `json:"roleId,omitempty"`
	DateCreated  time.Time     `json:"dateCreated"`
	DateModified time.Time     `json:"dateModified"`
}

type UserRepo struct {
	DB *pgxpool.Pool
}

func (r *UserRepo) FindByUsername(username string) (User, error) {
	sql := "select * from user_account where username = $1"

	var u User
	err := r.DB.QueryRow(context.Background(), sql, username).
		Scan(&u.Id, &u.Username, &u.Email, &u.RoleId, &u.DateCreated, &u.DateModified)

	fmt.Println(u)

	return u, err
}

func (r *UserRepo) GetPasswordForUser(username string) (string, error) {
	sql := "select password from user_account where username = $1"
	var password string
	err := r.DB.QueryRow(context.Background(), sql, username).Scan(&password)
	return password, err
}

func (r *UserRepo) Create(user User, hashedPassword string) (User, error) {
	sql := "insert into user_account (username, password, email, role_id, date_created, date_modified) " +
		"VALUES($1, $2, $3, $4, $5, $6) " +
		"RETURNING id, username, email, role_id, date_created, date_modified"

	var u User
	err := r.DB.QueryRow(context.Background(), sql, user.Username, hashedPassword, user.Email, user.RoleId, time.Now(), time.Now()).
		Scan(&u.Id, &u.Username, &u.Email, &u.RoleId, &u.DateCreated, &u.DateModified)

	fmt.Println(u)

	return u, err
}
