package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/Flynn81/tkdo/handler"
	"github.com/Flynn81/tkdo/model"
	"github.com/gorilla/handlers"

	_ "github.com/lib/pq"
)

func closeDB(db *sql.DB) {
	log.Println("closing db")
	err := db.Close()
	if err != nil {
		log.Println(err)
	}
}

const (
	envHost     = "TKDO_HOST"
	envPort     = "TKDO_PORT"
	envUser     = "TKDO_USER"
	envPassword = "TKDO_PASSWORD"
	envDbName   = "TKDO_DBNAME"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Server startup")

	// - TKDO_HOST
	// - TKDO_PORT
	// - TKDO_USER
	// - TKDO_PASSWORD
	// - TKDO_DBNAME

	//TODO: add some validation
	host := os.Getenv(envHost)
	port := os.Getenv(envPort)
	user := os.Getenv(envUser)
	password := os.Getenv(envPassword)
	dbname := os.Getenv(envDbName)

	var err2 error

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err2 := sql.Open("postgres",
		psqlInfo)
	if err2 != nil {
		log.Fatal("error connecting to the database: ", err2)
	}
	model.Init(db)
	defer closeDB(db)

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 3 && argsWithoutProg[0] == "admin" {
		initAdmin(argsWithoutProg[1], argsWithoutProg[2])
		return
	}

	ta := model.CockroachTaskAccess{}
	ua := model.CockroachUserAccess{}

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

	v := reflect.ValueOf(http.DefaultServeMux).Elem()
	fmt.Printf("routes: %v\n", v.FieldByName("m"))

	err := http.ListenAndServe(":7056", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)) //todo: make port env var
	if err != nil {
		panic(err)
	}
	log.Println("Server shutdown")
}

func initAdmin(p string, e string) {
	ua := model.CockroachUserAccess{}
	u := model.User{ID: "", Name: "admin", Email: e, Status: "admin"}
	c := ua.Create(&u)
	ua.UpdatePassword(c.ID, p)
}
