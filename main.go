package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Flynn81/tkdo/handler"
	"github.com/Flynn81/tkdo/model"
)

func main() {
	log.Println("Server startup")

	inMem := model.InMemTaskAccess{Mux: &sync.Mutex{}}

	lh := handler.ListHandler{TaskAccess: inMem}
	http.Handle("/tasks", handler.CheckMethod(handler.CheckHeaders(lh), "GET", "POST"))

	th := handler.TaskHandler{TaskAccess: inMem}
	http.Handle("/tasks/", handler.CheckMethod(handler.CheckHeaders(th), "GET", "DELETE", "PUT"))

	sh := handler.SearchHandler{TaskAccess: inMem}
	http.Handle("/tasks/search", handler.CheckMethod(handler.CheckHeaders(sh), "GET"))

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(rw, "these aren't the droids you're looking for")
		if err != nil {
			panic(err)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
	log.Println("Server shutdown")
}
