package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Flynn81/tkdo/model"
)

//SearchHandler provides handler funcs for search
type SearchHandler struct {
	TaskAccess model.TaskAccess
}

func (sh SearchHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	n := r.URL.Query().Get("name")
	t := r.URL.Query().Get("type")

	//TODO pagination

	tasks := sh.TaskAccess.GetMany(n, t)

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")

	bytes, err := json.Marshal(tasks)

	if err != nil {
		http.Error(rw, "error with search", http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(bytes)
	if err != nil {
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

	id, err := strconv.Atoi(segments[len(segments)-1])

	if err != nil {
		http.Error(rw, "error with get", http.StatusBadRequest)
		return
	}

	t, err := ta.Get(id)

	if err != nil {
		http.Error(rw, "error with get", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(t)

	if err != nil {
		http.Error(rw, "error with get", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)

	_, err = rw.Write(bytes)

	if err != nil {
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

	if t.ID == 0 {
		segments := strings.Split(r.RequestURI, "/")
		id, err2 := strconv.Atoi(segments[len(segments)-1])
		if err2 != nil {
			http.Error(rw, "error with update", http.StatusInternalServerError)
			return
		}
		t.ID = id
	}

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
			panic(err)
		}
	} else {
		http.Error(rw, "error with update", http.StatusInternalServerError)
	}
}

func deleteTask(rw http.ResponseWriter, r *http.Request, ta model.TaskAccess) {
	segments := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(segments[len(segments)-1])
	if err != nil {
		http.Error(rw, "error with delete", http.StatusBadRequest)
		return
	}
	t, err := ta.Get(id)
	if err != nil {
		http.Error(rw, "error with delete", http.StatusBadRequest)
		return
	}
	ta.Delete(id)
	bytes, err := json.Marshal(t)
	if err != nil {
		http.Error(rw, "error with delete", http.StatusBadRequest)
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
	newTask := ta.Create(&t)
	var bytes []byte
	bytes, err = json.Marshal(newTask)
	if err != nil {
		http.Error(rw, "error with create", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusCreated)
	_, err = rw.Write(bytes)
	if err != nil {
		panic(err)
	}
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

	tasks := ta.List(page, pageSize)

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
