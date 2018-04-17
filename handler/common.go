package handler

import "net/http"

func CheckMethod (method string, rw http.ResponseWriter, r *http.Request) error {
	if r.Method != method {

	}
	return nil
}
