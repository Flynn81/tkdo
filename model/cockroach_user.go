package model

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//CockroachUserAccess is a concrete struct implementing UserAccess, backed by CockroachDB
type CockroachUserAccess struct {
}

//Get returns a user given an email address
func (ua CockroachUserAccess) Get(email string) (*User, error) {
	rows, err := db.Query("select id, name, email, hash, status from task_user where email = $1", email)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer closeRows(rows)
	var (
		i, n, e, s string
	)
	var h []byte
	if rows.Next() {
		err := rows.Scan(&i, &n, &e, &h, &s)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		return &User{i, n, e, h, s}, nil
	}
	return nil, fmt.Errorf("user with email %v cannot be found", email)
}

//Create takes a user without an id and persists it
func (ua CockroachUserAccess) Create(u *User) *User {
	stmt, err := db.Prepare("INSERT INTO TASK_USER (id, name, email, status) VALUES ($1, $2, $3, $4)")

	if err != nil {
		log.Print(err)
		return nil
	}
	_, err = stmt.Exec(uuid.NewString(), u.Name, u.Email, u.Status)
	if err != nil {
		log.Print(err)
		return nil
	}
	createdUser, err := ua.Get(u.Email)
	if err != nil {
		log.Print(err)
		return nil
	}
	return createdUser
}

//UpdatePassword takes a task and attempt to update it
func (ua CockroachUserAccess) UpdatePassword(i string, p string) bool {
	stmt, err := db.Prepare("UPDATE TASK_USER SET hash=$1 WHERE id = $2")
	if err != nil {
		log.Print(err)
		return false
	}
	h, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	if err != nil {
		log.Print(err)
		return false
	}
	_, err = stmt.Exec(h, i)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}
