package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Flynn81/tkdo/model"
)

type SearchHandler struct {
	TaskAccess model.TaskAccess
}

func (sh SearchHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	n := r.URL.Query().Get("name")
	t := r.URL.Query().Get("type")

	//TODO pagination

	tasks := sh.TaskAccess.GetMany(n, t)

	bytes, err := json.Marshal(tasks)

	if err != nil {
		fmt.Fprintf(rw, "error with search: %v", err)
	}

	rw.Write(bytes)
}

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
	fmt.Printf("%q\n", segments)
	id, err := strconv.Atoi(segments[len(segments)-1])
	if err != nil {
		fmt.Fprintf(rw, "error with get: %v", err)
	}
	t, err := ta.Get(id)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rw, "error with get: %v", err)
	}
	bytes, err := json.Marshal(t)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "error with get: %v", err)
	}
	rw.Write(bytes)
}

func updateTask(rw http.ResponseWriter, r *http.Request, ta model.TaskAccess) {
	var t model.Task
	json.NewDecoder(r.Body).Decode(&t)

	if t.Id == 0 {
		segments := strings.Split(r.RequestURI, "/")
		fmt.Printf("%q\n", segments)
		id, err := strconv.Atoi(segments[len(segments)-1])
		if err != nil {
			fmt.Fprintf(rw, "error with get: %v", err)
		}
		t.Id = id
	}

	b := ta.Update(&t)
	bytes, err := json.Marshal(t)
	if err != nil {
		fmt.Fprintf(rw, "error with update: %v", err)
	}
	if b {
		rw.Write(bytes)
	} else {
		fmt.Fprintf(rw, "error with update: %v", err)
	}
}

func deleteTask(rw http.ResponseWriter, r *http.Request, ta model.TaskAccess) {
	segments := strings.Split(r.RequestURI, "/")
	fmt.Printf("%q\n", segments)
	id, err := strconv.Atoi(segments[len(segments)-1])
	if err != nil {
		fmt.Fprintf(rw, "error with get: %v", err)
	}
	if err != nil {
		fmt.Fprintf(rw, "error with delete: %v", err)
	}
	t, err := ta.Get(id)
	if err != nil {
		fmt.Fprintf(rw, "error with delete: %v", err)
	}
	ta.Delete(id)
	bytes, err := json.Marshal(t)
	if err != nil {
		fmt.Fprintf(rw, "error with delete: %v", err)
	}
	rw.Write(bytes)
}

func createTask(rw http.ResponseWriter, r *http.Request, ta model.TaskAccess) {
	var t model.Task
	json.NewDecoder(r.Body).Decode(&t)
	newTask := ta.Create(&t)
	bytes, err := json.Marshal(newTask)
	if err != nil {
		fmt.Fprintf(rw, "error with create: %v", err)
	}
	rw.Write(bytes)
}

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
		fmt.Fprintf(rw, "error with search: %v", err)
	}

	rw.Write(bytes)
}
