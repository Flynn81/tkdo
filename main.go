package main

import (
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/felixge/httpsnoop"
	"go.uber.org/zap"

	"github.com/Flynn81/tkdo/handler"
	"github.com/Flynn81/tkdo/model"
)

const (
	envHost       = "TKDO_HOST"
	envPort       = "TKDO_PORT"
	envDbPort     = "TKDO_DB_PORT"
	envUser       = "TKDO_USER"
	envDbName     = "TKDO_DBNAME"
	envDynamoHost = "TKDO_DYNAMOHOST"
	defaultPort   = "7056"
)

func main() {
	logger, _ := zap.NewProduction()
	//defer logger.Sync() // flushes buffer, if any

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	zap.S().Info("Server startup")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

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
	dbname := os.Getenv(envDbName)
	dynamoHost := os.Getenv(envDynamoHost)
	zap.S().Infof("%s:%s", envHost, host)
	zap.S().Infof("%s:%s", envPort, port)
	zap.S().Infof("%s:%s", envDbPort, dbPort)
	zap.S().Infof("%s:%s", envUser, user)
	zap.S().Infof("%s:%s", envDbName, dbname)
	zap.S().Infof("%s:%s", envDynamoHost, dynamoHost)

	if dynamoHost == "" {
		zap.S().Error("TKDO_DYNAMOHOST not set")
		panic("TKDO_DYNAMOHOST not set")
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess, &aws.Config{Endpoint: aws.String(dynamoHost)})

	model.Init(svc)

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
	//TODO: fill this in or remove
	//return
	// ua := model.CockroachUserAccess{}
	// u := model.User{Name: "admin", Email: e, Status: "admin"}
	// c := ua.Create(&u)
	// ua.UpdatePassword(c.Email, p)
}
