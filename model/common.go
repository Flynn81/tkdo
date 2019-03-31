package model

import (
	"database/sql"
	"log"
)

//Error is a struct used for relaying error messages back to callers of the api
type Error struct {
	Msg string `json:"msg"`
}

var db *sql.DB

//Init is use to set the db for the package
func Init(d *sql.DB) {
	db = d
}

func closeRows(r *sql.Rows) {
	err := r.Close()
	if err != nil {
		log.Fatal(err)
	}
}
