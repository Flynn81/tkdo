package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/felixge/httpsnoop"
	"go.uber.org/zap"

	"github.com/Flynn81/tkdo/handler"
	"github.com/Flynn81/tkdo/model"

	_ "github.com/lib/pq"
)

func closeDB(db *sql.DB) {
	zap.S().Info("closing db")
	err := db.Close()
	if err != nil {
		zap.S().Infof("%e", err)
	}
}

const (
	envHost     = "TKDO_HOST"
	envPort     = "TKDO_PORT"
	envDbPort   = "TKDO_DB_PORT"
	envUser     = "TKDO_USER"
	envPassword = "TKDO_PASSWORD"
	envDbName   = "TKDO_DBNAME"
	defaultPort = "7056"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	zap.S().Info("Server startup")

	// - TKDO_HOST
	// - TKDO_PORT
	// - TKDO_USER
	// - TKDO_PASSWORD
	// - TKDO_DBNAME

	host := os.Getenv(envHost)
	port := os.Getenv(envPort)
	if len(port) == 0 {
		port = defaultPort
	}
	dbPort := os.Getenv(envDbPort)
	user := os.Getenv(envUser)
	password := os.Getenv(envPassword)
	dbname := os.Getenv(envDbName)
	zap.S().Infof("%s:%s", envHost, host)
	zap.S().Infof("%s:%s", envPort, port)
	zap.S().Infof("%s:%s", envDbPort, dbPort)
	zap.S().Infof("%s:%s", envUser, user)
	zap.S().Infof("%s:%s", envDbName, dbname)

	var err2 error

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, dbPort, user, password, dbname)

	db, err2 := sql.Open("postgres",
		psqlInfo)
	if err2 != nil {
		zap.S().Infof("error connecting to the database: %s", err2)
		panic(err2)
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
			zap.S().Infof("We are panicked: %e", err)
			panic(err)
		}
	})

	v := reflect.ValueOf(http.DefaultServeMux).Elem()

	zap.S().Infof("routes: %v\n", v.FieldByName("m"))

	wrappedH := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(http.DefaultServeMux, w, r)
		zap.S().Infof(
			"%s %s (code=%d dt=%s written=%d)",
			r.Method,
			r.URL,
			m.Code,
			m.Duration,
			m.Written,
		)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), wrappedH) //todo: make port env var
	if err != nil {
		zap.S().Infof("We are panicked: %e", err)
		panic(err)
	}
	zap.S().Info("Server shutdown")
}

func initAdmin(p string, e string) {
	ua := model.CockroachUserAccess{}
	u := model.User{ID: "", Name: "admin", Email: e, Status: "admin"}
	c := ua.Create(&u)
	ua.UpdatePassword(c.ID, p)
}
