package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/Flynn81/tkdo/model"
)

//CheckCors will check if cors flag is true and if so add CORS headers in the response
func CheckCors(h http.Handler, cors bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cors {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		h.ServeHTTP(w, r)
	})
}

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
			zap.S().Info("content-type not accepted")
			http.Error(w, "content-type not accepted", http.StatusBadRequest)
			return
		} else if !u && r.Header.Get("uid") == "" { //todo look up valid uids from an in memory cache
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			bytes, err := json.Marshal(model.Error{Msg: "invalid uid"})
			if err != nil {
				zap.S().Infof("We are panicked: %e", err)
				panic(err)
			}
			_, err = w.Write(bytes)
			if err != nil {
				zap.S().Infof("We are panicked: %e", err)
				panic(err)
			}
			return
		}

		//write error response
		//set response status code

		h.ServeHTTP(w, r)
	})
}
