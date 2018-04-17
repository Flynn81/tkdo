package handler

import (
	"net/http"
)

//CheckMethod takes a slice of strings, which should represent the http
//request methods that the handler will allow.
func CheckMethod(h http.Handler, methods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, method := range methods {
			if r.Method == method {
				h.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})
}

//CheckHeaders ensures all requests being handled have the correct headers
func CheckHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json; charset=utf-8" {
			http.Error(w, "content-type not accepted", http.StatusBadRequest)
		}

		//write error response
		//set response status code

		h.ServeHTTP(w, r)
	})
}
