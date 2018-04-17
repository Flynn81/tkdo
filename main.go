package main

import (
	"net/http"
	"fmt"
	"log"
	"tkdo.handler"
)

func main() {
	log.Println("Server startup")
	http.HandleFunc("/tasks", func(rw http.ResponseWriter, r *http.Request) {
		handler.CheckMethod(http.MethodGet, rw, r)
		//CHECK headers
		//CHECK REQUEST METHOD ALLOWED
		fmt.Fprintf(rw, "%v:tasks", r.Method)
	})
	http.HandleFunc("/tasks/search", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "%v:tasks/search", r.Method)
	})
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "Hello")
	})
	http.ListenAndServe(":8080", nil)
	log.Println("Server shutdown")
}