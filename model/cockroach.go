package model

import (
	"fmt"
	"log"

	"database/sql"
)

//CockroachTaskAccess is a concrete strut implementing TaskAccess, backed by CockroachDB
type CockroachTaskAccess struct {
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

//Get returns an task given an id.
func (ta CockroachTaskAccess) Get(id string) (*Task, error) {

	rows, err := db.Query("select id, name, type from task where id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer closeRows(rows)

	var (
		i, n, t string
	)
	if rows.Next() {
		err := rows.Scan(&i, &n, &t)
		if err != nil {
			log.Fatal(err)
		}
		return &Task{i, n, t}, nil
	}

	return nil, fmt.Errorf("task with id %v cannot be found", id)
}

//Create takes a task without id and persists it.
func (ta CockroachTaskAccess) Create(t *Task) *Task {
	log.Println(t)
	stmt, err := db.Prepare("INSERT INTO TASK (name, type) VALUES ($1, $2)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(t.Name, t.TaskType)
	if err != nil {
		log.Fatal(err)
	}

	return t
}

//Update takes a task and attempt to update it
func (ta CockroachTaskAccess) Update(task *Task) bool {
	stmt, err := db.Prepare("UPDATE TASK SET name=$1, type=$2 WHERE id = $3")
	if err != nil {
		log.Fatal(err)
		return false
	}
	_, err = stmt.Exec(task.Name, task.TaskType, task.ID)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

//Delete takes an id and attempts to delete the task with the id
func (ta CockroachTaskAccess) Delete(id string) bool {
	stmt, err := db.Prepare("DELETE FROM TASK WHERE id = $1")
	if err != nil {
		log.Fatal(err)
		return false
	}
	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

//GetMany returns tasks matching the given string and/or task type
func (ta CockroachTaskAccess) GetMany(keyword string, taskType string) []*Task {

	q := "select id, name, type from task"
	if keyword == "" && taskType == "" {
		return nil
	}
	q = q + " where "
	if keyword != "" {
		q = q + " name LIKE '%" + keyword + "%'"
	}
	if taskType != "" {
		q = q + " and type = '" + taskType + "'"
	}
	rows, err := db.Query(q)
	fmt.Println("keyword is ", keyword)
	fmt.Println("q is ", q)
	var r = []*Task{}
	if err != nil {
		log.Fatal(err)
	}
	defer closeRows(rows)
	var (
		i, n, t string
	)
	for rows.Next() {
		err := rows.Scan(&i, &n, &t)
		if err != nil {
			log.Fatal(err)
		}
		r = append(r, &Task{i, n, t})
	}
	return r
}

//List returns a list of tasks
func (ta CockroachTaskAccess) List(page int, pageSize int) []*Task {
	log.Printf("task length is %v", len(tasks))
	log.Printf("page size is %v", pageSize)
	log.Printf("page is %v", page)
	m := (page * pageSize) + pageSize
	if m > len(tasks) {
		m = len(tasks)
	}
	log.Print("m is ", m)

	rows, err := db.Query("select id, name, type from task")
	var r = []*Task{}
	if err != nil {
		log.Fatal(err)
	}
	defer closeRows(rows)
	var (
		i, n, t string
	)
	rownum := 0
	for rows.Next() {
		err := rows.Scan(&i, &n, &t)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("id is ", i)
		if rownum >= page*pageSize && rownum < page*pageSize+pageSize {
			r = append(r, &Task{i, n, t})
		}
		rownum++
	}
	return r
}
