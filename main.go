package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Flynn81/tkdo/handler"
	"github.com/Flynn81/tkdo/model"

	_ "github.com/lib/pq"
)

func closeDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Server startup")

	var err2 error
	db, err2 := sql.Open("postgres",
		"postgresql://tk@localhost:26257/tkdo?sslmode=disable")
	if err2 != nil {
		log.Fatal("error connecting to the database: ", err2)
	}
	model.Init(db)
	defer closeDB(db)

	ta := model.CockroachTaskAccess{}
	ua := model.CockroachUserAccess{}
	//inMem := model.InMemTaskAccess{Mux: &sync.Mutex{}}

	lh := handler.ListHandler{TaskAccess: ta}
	http.Handle("/tasks", handler.CheckMethod(handler.CheckHeaders(lh, false), "GET", "POST"))

	th := handler.TaskHandler{TaskAccess: ta}
	http.Handle("/tasks/", handler.CheckMethod(handler.CheckHeaders(th, false), "GET", "DELETE", "PUT"))

	sh := handler.SearchHandler{TaskAccess: ta}
	http.Handle("/tasks/search", handler.CheckMethod(handler.CheckHeaders(sh, false), "GET"))

	uh := handler.UserHandler{UserAccess: ua}
	http.Handle("/users", handler.CheckMethod(handler.CheckHeaders(uh, true), "POST"))

	http.HandleFunc("/hc", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
		_, err := fmt.Fprint(rw, "these aren't the droids you're looking for")
		if err != nil {
			panic(err)
		}
	})

	err := http.ListenAndServe(":7056", nil)
	if err != nil {
		panic(err)
	}
	log.Println("Server shutdown")
}
