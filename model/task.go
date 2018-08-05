package model

import (
	"fmt"
	"strings"
)

var tasks []*Task

type Task struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	TaskType string `json:"taskType"`
}

type TaskAccess interface {
	Get(id int) (*Task, error)
	Create(t *Task) *Task
	Update(task *Task) bool
	Delete(id int) bool
	GetMany(keyword string, taskType string) []*Task
	List(page int, pageSize int) []*Task
}

type InMemTaskAccess struct {
}

func (ta InMemTaskAccess) Get(id int) (*Task, error) {
	for _, t := range tasks {
		if t.Id == id {
			return t, nil
		}
	}
	return nil, fmt.Errorf("task with id %v cannot be found", id)
}

func (ta InMemTaskAccess) Create(t *Task) *Task {
	task := Task{len(tasks) + 1, t.Name, t.TaskType}
	fmt.Printf("task length is %v", len(tasks))
	tasks = append(tasks, &task)
	fmt.Printf("task length is %v", len(tasks))
	return &task
}

func (ta InMemTaskAccess) Update(task *Task) bool {
	for _, t := range tasks {
		if t.Id == task.Id {
			t.Name = task.Name
			t.TaskType = task.TaskType
			return true
		}
	}
	return false
}

func (ta InMemTaskAccess) Delete(id int) bool {
	for i, t := range tasks {
		if t.Id == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return true
		}
	}
	return false
}

func (ta InMemTaskAccess) GetMany(keyword string, taskType string) []*Task {
	fmt.Println("keyword is ", keyword)
	var r []*Task
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

func (ta InMemTaskAccess) List(page int, pageSize int) []*Task {
	fmt.Println("task length is %v", len(tasks))
	fmt.Println("page size is %v", pageSize)
	fmt.Println("page is %v", page)
	m := (page * pageSize) + pageSize
	if m > len(tasks) {
		m = len(tasks)
	}
	fmt.Println("m is ", m)
	if tasks != nil {
		return tasks[page*pageSize : m]
	}
	return nil
}
