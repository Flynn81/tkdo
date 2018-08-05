package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Flynn81/tkdo/handler"
	"github.com/Flynn81/tkdo/model"
)

func main() {
	log.Println("Server startup")

	//http.HandleFunc("/",CheckMethod(nil, "GET"))

	//MAKE INDIVIDUAL HANDLERS, CHAIN THE CALLS TO HANDLE

	lh := handler.ListHandler{model.InMemTaskAccess{}}
	http.Handle("/tasks", handler.CheckMethod(handler.CheckHeaders(lh), "GET", "POST"))

	th := handler.TaskHandler{model.InMemTaskAccess{}}
	http.Handle("/tasks/", handler.CheckMethod(handler.CheckHeaders(th), "GET", "DELETE", "PUT"))

	sh := handler.SearchHandler{model.InMemTaskAccess{}}
	http.Handle("/tasks/search", handler.CheckMethod(handler.CheckHeaders(sh), "GET"))

	//http.HandleFunc("/tasks/search", sh.Search)
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "these aren't the droids you're looking for")
	})

	http.ListenAndServe(":8080", nil)
	log.Println("Server shutdown")
}
