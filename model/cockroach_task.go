package model

import (
	"fmt"
	"log"
)

//CockroachTaskAccess is a concrete strut implementing TaskAccess, backed by CockroachDB
type CockroachTaskAccess struct {
}

//Get returns an task given an id.
func (ta CockroachTaskAccess) Get(id string, userID string) (*Task, error) {

	rows, err := db.Query("select id, name, type, status, user_id from task where id = $1 and user_id = $2", id, userID)
	if err != nil {
		log.Fatal(err)
	}
	defer closeRows(rows)

	var (
		i, n, t, s, u string
	)
	if rows.Next() {
		err := rows.Scan(&i, &n, &t, &s, &u)
		if err != nil {
			log.Fatal(err)
		}
		return &Task{i, n, t, s, u}, nil
	}

	return nil, fmt.Errorf("task with id %v cannot be found", id)
}

//Create takes a task without id and persists it.
func (ta CockroachTaskAccess) Create(t *Task) *Task {
	log.Println(t)
	stmt, err := db.Prepare("INSERT INTO TASK (name, type, status, user_id) VALUES ($1, $2, 'new', $3)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(t.Name, t.TaskType, t.UserID)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

//Update takes a task and attempt to update it
func (ta CockroachTaskAccess) Update(task *Task) bool {
	stmt, err := db.Prepare("UPDATE TASK SET name=$1, type=$2, status=$3 WHERE id = $4 and user_id = $5")
	if err != nil {
		log.Fatal(err)
		return false
	}
	_, err = stmt.Exec(task.Name, task.TaskType, task.Status, task.ID, task.UserID)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

//Delete takes an id and attempts to delete the task with the id
func (ta CockroachTaskAccess) Delete(id string, userID string) bool {
	stmt, err := db.Prepare("DELETE FROM TASK WHERE id = $1 and user_id = $2")
	if err != nil {
		log.Fatal(err)
		return false
	}
	_, err = stmt.Exec(id, userID)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

//GetMany returns tasks matching the given string and/or task type
func (ta CockroachTaskAccess) GetMany(keyword string, taskType string, userID string) []*Task {

	q := "select id, name, type, status, user_id from task"
	if keyword == "" && taskType == "" {
		return nil
	}
	q = q + " where user_id = $1 "
	if keyword != "" {
		q = q + " and name LIKE '%" + keyword + "%'"
	}
	if taskType != "" {
		q = q + " and type = '" + taskType + "'"
	}
	rows, err := db.Query(q, userID)
	fmt.Println("keyword is ", keyword)
	fmt.Println("q is ", q)
	var r = []*Task{}
	if err != nil {
		log.Fatal(err)
	}
	defer closeRows(rows)
	var (
		i, n, t, s, u string
	)
	for rows.Next() {
		err := rows.Scan(&i, &n, &t, &s, &u)
		if err != nil {
			log.Fatal(err)
		}
		r = append(r, &Task{i, n, t, s, u})
	}
	return r
}

//List returns a list of tasks
func (ta CockroachTaskAccess) List(page int, pageSize int, userID string) []*Task {
	log.Printf("task length is %v", len(tasks))
	log.Printf("page size is %v", pageSize)
	log.Printf("page is %v", page)
	m := (page * pageSize) + pageSize
	if m > len(tasks) {
		m = len(tasks)
	}
	log.Print("m is ", m)

	rows, err := db.Query("select id, name, type, status, user_id from task where user_id = $1", userID)
	var r = []*Task{}
	if err != nil {
		log.Fatal(err)
	}
	defer closeRows(rows)
	var (
		i, n, t, s, u string
	)
	rownum := 0
	for rows.Next() {
		err := rows.Scan(&i, &n, &t, &s, &u)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("id is ", i)
		if rownum >= page*pageSize && rownum < page*pageSize+pageSize {
			r = append(r, &Task{i, n, t, s, u})
		}
		rownum++
	}
	return r
}
