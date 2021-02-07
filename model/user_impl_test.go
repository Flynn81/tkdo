package model

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateUser(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	name := "alice"
	email := "any@email.com"
	status := "new"
	h := []byte{}
	mock.ExpectPrepare("INSERT INTO TASK_USER").ExpectExec().WithArgs(sqlmock.AnyArg(), name, email, status).WillReturnResult(sqlmock.NewResult(0, 1))
	columns := []string{"id", "name", "email", "hash", "status"}
	mock.ExpectQuery("select id, name, email, hash, status from task_user").WithArgs(email).WillReturnRows(sqlmock.NewRows(columns).AddRow("an id", name, email, h, status))

	ua := CockroachUserAccess{}
	user := User{Name: name, Email: email, Status: status}
	returnedUser := ua.Create(&user)
	if returnedUser.Email != email {
		t.Errorf("email not set correctly")
	} else if returnedUser.Name != name {
		t.Errorf("name not set correctly")
	} else if returnedUser.Status != status {
		t.Errorf("status not set correctly")
	} else if returnedUser.ID == "" {
		t.Errorf("ID not set correctly")
	}
}

func TestGetUser(t *testing.T) {
	var mock sqlmock.Sqlmock
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	name := "alice"
	email := "any@email.com"
	status := "new"
	h := []byte{}
	columns := []string{"id", "name", "email", "hash", "status"}
	mock.ExpectQuery("select id, name, email, hash, status from task_user").WithArgs(email).WillReturnRows(sqlmock.NewRows(columns).AddRow("an id", name, email, h, status))

	ua := CockroachUserAccess{}
	returnedUser, err := ua.Get(email)
	if err != nil {
		t.Errorf("error when getting user, %e", err)
	} else if returnedUser.Email != email {
		t.Errorf("email not set correctly")
	} else if returnedUser.Name != name {
		t.Errorf("name not set correctly")
	} else if returnedUser.Status != status {
		t.Errorf("status not set correctly")
	} else if returnedUser.ID == "" {
		t.Errorf("ID not set correctly")
	}
}
