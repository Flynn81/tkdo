package model

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

var tasks []*Task

//Task represents the data structure for a task
type Task struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	TaskType string `json:"taskType"`
}

//TaskAccess functions to work with tasks
type TaskAccess interface {
	Get(id int) (*Task, error)
	Create(t *Task) *Task
	Update(task *Task) bool
	Delete(id int) bool
	GetMany(keyword string, taskType string) []*Task
	List(page int, pageSize int) []*Task
}

//InMemTaskAccess is a concrete struct implementing TaskAccess using
//an in-memory data store in the form of an array.
type InMemTaskAccess struct {
	Mux *sync.Mutex
}

//Get returns an task given an id.
func (ta InMemTaskAccess) Get(id int) (*Task, error) {
	for _, t := range tasks {
		if t.ID == id {
			fmt.Println("task found")
			return t, nil
		}
	}
	return nil, fmt.Errorf("task with id %v cannot be found", id)
}

//Create takes a task without id and persists it.
func (ta InMemTaskAccess) Create(t *Task) *Task {
	ta.Mux.Lock()
	task := Task{len(tasks) + 1, t.Name, t.TaskType}
	tasks = append(tasks, &task)
	ta.Mux.Unlock()
	return &task
}

//Update takes a task and attempt to update it
func (ta InMemTaskAccess) Update(task *Task) bool {
	ta.Mux.Lock()
	for _, t := range tasks {
		if t.ID == task.ID {
			t.Name = task.Name
			t.TaskType = task.TaskType
			ta.Mux.Unlock()
			return true
		}
	}
	ta.Mux.Unlock()
	return false
}

//Delete takes an id and attempts to delete the task with the id
func (ta InMemTaskAccess) Delete(id int) bool {
	ta.Mux.Lock()
	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			ta.Mux.Unlock()
			return true
		}
	}
	ta.Mux.Unlock()
	return false
}

//GetMany returns tasks matching the given string and/or task type
func (ta InMemTaskAccess) GetMany(keyword string, taskType string) []*Task {
	fmt.Println("keyword is ", keyword)
	var r = []*Task{}
	for _, t := range tasks {
		var k, tt = false, false //todo init before for loop
		if keyword != "" && strings.Contains(t.Name, keyword) {
			k = true
		} else if keyword == "" { //todo do this before the loop
			k = true
		}
		if taskType != "" && t.TaskType == taskType {
			tt = true
		} else if taskType == "" {
			tt = true
		}
		if k && tt {
			r = append(r, t)
		}
	}
	return r
}

//List returns a list of tasks
func (ta InMemTaskAccess) List(page int, pageSize int) []*Task {
	log.Printf("task length is %v", len(tasks))
	log.Printf("page size is %v", pageSize)
	log.Printf("page is %v", page)
	m := (page * pageSize) + pageSize
	if m > len(tasks) {
		m = len(tasks)
	}
	log.Print("m is ", m)
	if tasks != nil {
		return tasks[page*pageSize : m]
	}
	log.Print("empty array")
	return []*Task{}
}
