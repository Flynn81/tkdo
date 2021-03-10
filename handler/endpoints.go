package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Flynn81/tkdo/model"
	"go.uber.org/zap"
)

const (
	regularTaskType = "regular"
)

//UserHandler processes http requests for user operations
type UserHandler struct {
	UserAccess model.UserAccess
}

func (uh UserHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var u model.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		zap.S().Infof("error with create decoding the request body, %e", err)
		http.Error(rw, "error with create", http.StatusInternalServerError)
		return
	}

	newUser := uh.UserAccess.Create(&u)

	var bytes []byte
	bytes, err = json.Marshal(newUser)
	if err != nil {
		zap.S().Infof("error with create, %e", err)
		http.Error(rw, "error with create", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusCreated)
	_, err = rw.Write(bytes)
	if err != nil {
		zap.S().Infof("We are panicked: %e", err)
		panic(err)
	}
}

//SearchHandler provides handler funcs for search
type SearchHandler struct {
	TaskAccess model.TaskAccess
}

func (sh SearchHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	n := r.URL.Query().Get("name")
	t := r.URL.Query().Get("type")

	//TODO pagination

	if n[0] == "'"[0] {
		n = n[1 : len(n)-1]
	}

	tasks := sh.TaskAccess.GetMany(n, t, r.Header.Get("uid"))

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")

	bytes, err := json.Marshal(tasks)

	if err != nil {
		http.Error(rw, "error with search", http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(bytes)
	if err != nil {
		zap.S().Infof("We are panicked: %e", err)
		panic(err)
	}
}

//TaskHandler provides handler funcs for tasks
type TaskHandler struct {
	TaskAccess model.TaskAccess
}

func (th TaskHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { //is there a const I can use for this?
		getTask(rw, r, th.TaskAccess)
	} else if r.Method == "DELETE" {
		deleteTask(rw, r, th.TaskAccess)
	} else if r.Method == "PUT" {
		updateTask(rw, r, th.TaskAccess)
	}
}

func getTask(rw http.ResponseWriter, r *http.Request, ta model.TaskAccess) {

	segments := strings.Split(r.RequestURI, "/")

	id := strings.Split(segments[len(segments)-1], "?")[0]

	t, err := ta.Get(id, r.Header.Get("uid"))

	if err != nil {
		http.Error(rw, "error with get: "+err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(t)

	if err != nil {
		http.Error(rw, "error with get: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)

	_, err = rw.Write(bytes)

	if err != nil {
		zap.S().Infof("We are panicked: %e", err)
		panic(err)
	}
}

func updateTask(rw http.ResponseWriter, r *http.Request, ta model.TaskAccess) {
	var t model.Task
	err := json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		http.Error(rw, "error with update", http.StatusInternalServerError)
		return
	}

	if t.ID == "" {
		segments := strings.Split(r.RequestURI, "/")
		id := strings.Split(segments[len(segments)-1], "?")[0]
		t.ID = id
	}
	t.UserID = r.Header.Get("uid")
	b := ta.Update(&t)
	var bytes []byte
	bytes, err = json.Marshal(t)
	if err != nil {
		http.Error(rw, "error with update", http.StatusInternalServerError)
		return
	}
	if b {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err = rw.Write(bytes)
		if err != nil {
			zap.S().Infof("We are panicked: %e", err)
			panic(err)
		}
	} else {
		http.Error(rw, "error with update", http.StatusInternalServerError)
	}
}

func deleteTask(rw http.ResponseWriter, r *http.Request, ta model.TaskAccess) {
	segments := strings.Split(r.RequestURI, "/")
	id := strings.Split(segments[len(segments)-1], "?")[0]

	t, err := ta.Get(id, r.Header.Get("uid"))
	if err != nil {
		http.Error(rw, "error with delete: "+err.Error(), http.StatusBadRequest)
		return
	}
	ta.Delete(id, r.Header.Get("uid"))
	bytes, err := json.Marshal(t)
	if err != nil {
		http.Error(rw, "error with delete: "+err.Error(), http.StatusBadRequest)
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = rw.Write(bytes)
	if err != nil {
		panic(err)
	}
}

func createTask(rw http.ResponseWriter, r *http.Request, ta model.TaskAccess) {
	var t model.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(rw, "error with create", http.StatusInternalServerError)
		return
	}
	if t.TaskType == "" {
		t.TaskType = regularTaskType
	}
	t.UserID = r.Header.Get("uid")
	ta.Create(&t)
	// var bytes []byte
	// bytes, err = json.Marshal(newTask)
	// if err != nil {
	// 	http.Error(rw, "error with create", http.StatusInternalServerError)
	// 	return
	// }
	//rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusCreated)
	// _, err = rw.Write(bytes)
	// if err != nil {
	// 	panic(err)
	// }
}

//ListHandler is used to handle requests to list or create tasks
type ListHandler struct {
	TaskAccess model.TaskAccess
}

func (sh ListHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { //is there a const I can use for this?
		listTasks(rw, r, sh.TaskAccess)
	} else if r.Method == "POST" {
		createTask(rw, r, sh.TaskAccess)
	}
}

func listTasks(rw http.ResponseWriter, r *http.Request, ta model.TaskAccess) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 0
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		pageSize = 50
	}

	tasks := ta.List(page, pageSize, r.Header.Get("uid"))

	bytes, err := json.Marshal(tasks)

	if err != nil {
		http.Error(rw, "error with search", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")

	_, err = rw.Write(bytes)
	if err != nil {
		panic(err)
	}
}
