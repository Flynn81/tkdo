package model

import (
	"fmt"
	"log"
)

//CockroachUserAccess is a concrete struct implementing UserAccess, backed by CockroachDB
type CockroachUserAccess struct {
}

//Get returns a user given an email address
func (ua CockroachUserAccess) Get(email string) (*User, error) {
	rows, err := db.Query("select id, name, email from task_user where email = $1", email)
	if err != nil {
		log.Fatal(err)
	}
	defer closeRows(rows)
	var (
		i, n, e string
	)
	if rows.Next() {
		err := rows.Scan(&i, &n, &e)
		if err != nil {
			log.Fatal(err)
		}
		return &User{i, n, e}, nil
	}
	return nil, fmt.Errorf("user with email %v cannot be found", email)
}

//Create takes a user without an id and persists it
func (ua CockroachUserAccess) Create(u *User) *User {
	stmt, err := db.Prepare("INSERT INTO TASK_USER (name, email) VALUES ($1, $2) RETURNING ID")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(u.Name, u.Email)
	if err != nil {
		log.Fatal(err)
	}
	createdUser, err := ua.Get(u.Email)
	if err != nil {
		log.Fatal(err)
	}
	return createdUser
}
