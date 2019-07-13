package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Flynn81/tkdo/model"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/oauth2.v3/server"
)

//LoginHandler processes http requests for user operations
type LoginHandler struct {
	UserAccess model.UserAccess
	Oauth      *server.Server
}

func (h LoginHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var login model.Login
	err := json.NewDecoder(r.Body).Decode(&login)

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err != nil {
		http.Error(rw, "error with decoding login request body", http.StatusInternalServerError)
		return
	}

	u, err := h.UserAccess.Get(login.Email)
	log.Printf("u is %v", u)

	if err != nil || u == nil {
		http.Error(rw, "error with login", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword(u.Hash, []byte(login.Password))

	if err != nil {
		http.Error(rw, "error with login", http.StatusInternalServerError)
		return
	}

	//success now setup token and response
	q := r.URL.Query()
	q.Add("grant_type", "client_credentials")
	q.Add("client_id", "000000")
	q.Add("client_secret", "999999")
	r.URL.RawQuery = q.Encode()

	gt, tgr, verr := h.Oauth.ValidationTokenRequest(r)
	if verr != nil {
		m := fmt.Sprintf("error with login: %v", verr)
		http.Error(rw, m, http.StatusInternalServerError)
		return
	}

	ti, verr := h.Oauth.GetAccessToken(gt, tgr)
	if verr != nil {
		http.Error(rw, "error with login", http.StatusInternalServerError)
		return
	}

	res := model.LoginResponse{Name: u.Name, Email: u.Email, Token: ti.GetAccess()}

	bytes, err := json.Marshal(res)
	if err != nil {
		http.Error(rw, "error with login", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(bytes)
	if err != nil {
		panic(err)
	}
}
