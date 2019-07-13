package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Flynn81/tkdo/model"
	"gopkg.in/oauth2.v3/server"
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
func CheckHeaders(h http.Handler, u bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json; charset=utf-8" {
			http.Error(w, "content-type not accepted", http.StatusBadRequest)
		} else if !u && r.Header.Get("uid") == "" { //todo look up valid uids from an in memory cache
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			bytes, err := json.Marshal(model.Error{Msg: "invalid uid"})
			if err != nil {
				panic(err)
			}
			_, err = w.Write(bytes)
			if err != nil {
				panic(err)
			}
			return
		}

		//write error response
		//set response status code

		h.ServeHTTP(w, r)
	})
}

//CheckForAuthToken is middleware to check for an auth token
func CheckForAuthToken(h http.Handler, s server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := s.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if int64(time.Until(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn())).Seconds()) < 0 {
			http.Error(w, "expired token", http.StatusBadRequest)
			return
		}
		h.ServeHTTP(w, r)
	})
}
