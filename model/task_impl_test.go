package model

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateTask(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	name := "alice"
	taskType := "basic"
	status := "new"
	userID := "test-user-id"

	mock.ExpectPrepare("INSERT INTO TASK").ExpectExec().WithArgs(sqlmock.AnyArg(), name, taskType, userID).WillReturnResult(sqlmock.NewResult(0, 1))

	ta := CockroachTaskAccess{}
	task := Task{Name: name, TaskType: taskType, Status: status, UserID: userID}
	returnedTask := ta.Create(&task)

	if returnedTask.Name != name {
		t.Errorf("name not set correctly")
	} else if returnedTask.TaskType != taskType {
		t.Errorf("task type not set correctly")
	} else if returnedTask.Status != status {
		t.Errorf("status not set correctly")
	} else if returnedTask.UserID != userID {
		t.Errorf("user id not set correctly")
	} else if returnedTask.ID == "" {
		t.Errorf("ID not set correctly")
	}
}

func TestGetTask(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	name := "alice"
	taskType := "basic"
	status := "new"
	userID := "test-user-id"
	id := "an-id"

	columns := []string{"id", "name", "email", "hash", "status"}
	mock.ExpectQuery("select id, name, type, status, user_id from task").WithArgs(id, userID).WillReturnRows(sqlmock.NewRows(columns).AddRow(id, name, taskType, status, userID))

	ta := CockroachTaskAccess{}
	returnedTask, err := ta.Get(id, userID)

	if err != nil {
		t.Errorf("error getting task, %e", err)
	} else if returnedTask.Name != name {
		t.Errorf("name not set correctly")
	} else if returnedTask.TaskType != taskType {
		t.Errorf("task type not set correctly")
	} else if returnedTask.Status != status {
		t.Errorf("status not set correctly")
	} else if returnedTask.UserID != userID {
		t.Errorf("user id not set correctly")
	} else if returnedTask.ID == "" {
		t.Errorf("ID not set correctly")
	}
}

func TestUpdateTask(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	name := "alice"
	taskType := "basic"
	status := "new"
	userID := "test-user-id"
	id := "an-id"
	mock.ExpectPrepare("UPDATE TASK SET name").ExpectExec().WithArgs(name, taskType, status, id, userID).WillReturnResult(sqlmock.NewResult(0, 1))

	task := Task{ID: id, Name: name, TaskType: taskType, Status: status, UserID: userID}
	ta := CockroachTaskAccess{}
	if !ta.Update(&task) {
		t.Errorf("update returned false")
	}
}

func TestDeleteTask(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := "test-user-id"
	id := "an-id"
	mock.ExpectPrepare("DELETE FROM TASK").ExpectExec().WithArgs(id, userID).WillReturnResult(sqlmock.NewResult(0, 1))

	ta := CockroachTaskAccess{}

	if !ta.Delete(id, userID) {
		t.Errorf("delete returned false")
	}
}
