package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type User struct {
	Id           uuid.UUID      `json:"id"`
	Username     string         `json:"username" binding:"required"`
	Password     string         `json:"password,omitempty"`
	Email        sql.NullString `json:"email" binding:"required"`
	RoleId       uuid.UUID      `json:"roleId,omitempty" binding:"required"`
	DateCreated  time.Time      `json:"dateCreated"`
	DateModified time.Time      `json:"dateModified"`
}

type UserRepo struct {
	DB *pgxpool.Pool
}

func (r *UserRepo) FindAll() ([]User, error) {
	userRows, err := r.DB.Query(context.Background(), "select * from user_account")
	if err != nil {
		return nil, err
	}
	defer userRows.Close()

	var users []User
	for userRows.Next() {
		u := new(User)
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

func (r *UserRepo) FindByUsername(username string) (User, error) {
	ds := "select * from user_account where username = $1"

	var u User
	err := r.DB.QueryRow(context.Background(), ds, username).
		Scan(&u.Id, &u.Username, &u.Email, &u.RoleId, &u.DateCreated, &u.DateModified)

	fmt.Println(u)

	return u, err
}

func (r *UserRepo) GetUserLoginInfo(username string) (User, error) {
	ds := "select id, role_id, password from user_account where username = $1"
	var user User
	err := r.DB.QueryRow(context.Background(), ds, username).Scan(&user.Id, &user.RoleId, &user.Password)
	return user, err
}

// Create TODO(Handle email)
func (r *UserRepo) Create(user User, hashedPassword string) (User, error) {
	ds := "insert into user_account (username, password, email, role_id, date_created, date_modified) " +
		"VALUES($1, $2, $3, $4, $5, $6) " +
		"RETURNING id, username, email, role_id, date_created, date_modified"

	var u User
	err := r.DB.QueryRow(context.Background(), ds, user.Username, hashedPassword, user.Email, user.RoleId, time.Now(), time.Now()).
		Scan(&u.Id, &u.Username, &u.Email, &u.RoleId, &u.DateCreated, &u.DateModified)

	return u, err
}

func (u *User) EmptyPassword() {
	u.Password = ""
}
