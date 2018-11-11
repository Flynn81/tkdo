package model

var tasks []*Task

//Task represents the data structure for a task
type Task struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	TaskType string `json:"taskType"`
}

//TaskAccess functions to work with tasks
type TaskAccess interface {
	Get(id string) (*Task, error)
	Create(t *Task) *Task
	Update(task *Task) bool
	Delete(id string) bool
	GetMany(keyword string, taskType string) []*Task
	List(page int, pageSize int) []*Task
}
