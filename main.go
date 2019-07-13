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

	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"

	_ "github.com/lib/pq"
)

func closeDB(db *sql.DB) {
	log.Println("closing db")
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

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 3 && argsWithoutProg[0] == "admin" {
		initAdmin(argsWithoutProg[1], argsWithoutProg[2])
		return
	}

	oauth := initOauth()

	ta := model.CockroachTaskAccess{}
	ua := model.CockroachUserAccess{}

	lh := handler.ListHandler{TaskAccess: ta}
	http.Handle("/tasks", handler.CheckForAuthToken(handler.CheckMethod(handler.CheckHeaders(lh, false), "GET", "POST"), *oauth))

	th := handler.TaskHandler{TaskAccess: ta}
	http.Handle("/tasks/", handler.CheckForAuthToken(handler.CheckMethod(handler.CheckHeaders(th, false), "GET", "DELETE", "PUT"), *oauth))

	sh := handler.SearchHandler{TaskAccess: ta}
	http.Handle("/tasks/search", handler.CheckForAuthToken(handler.CheckMethod(handler.CheckHeaders(sh, false), "GET"), *oauth))

	uh := handler.UserHandler{UserAccess: ua}
	http.Handle("/users", handler.CheckForAuthToken(handler.CheckMethod(handler.CheckHeaders(uh, true), "POST"), *oauth))

	http.HandleFunc("/hc", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	logh := handler.LoginHandler{UserAccess: ua, Oauth: oauth}
	http.Handle("/login", handler.CheckMethod(handler.CheckHeaders(logh, true), "POST"))

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

func initOauth() *server.Server {

	log.Println("initOauth")
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	err := clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost:7056",
	})
	if err != nil {
		panic(err)
	}
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	//srv.SetUserAuthorizationHandler(userAuthorizeHandler)  I THINK THIS IS FOR LOGIN...

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	return srv
}
