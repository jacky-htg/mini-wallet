package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"mini-wallet/libraries/api"
	"mini-wallet/libraries/token"
	"mini-wallet/models"
	"mini-wallet/payloads/request"
	"mini-wallet/payloads/response"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Auths struct
type Auths struct {
	Db  *sql.DB
	Log *log.Logger
}

// Login http handler
func (u *Auths) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest request.LoginRequest
	err := loginRequest.Transform(r)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, err)
		return
	}

	uLogin := models.Customer{Username: loginRequest.Username}
	err = uLogin.GetByUsername(r.Context(), u.Db)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("call login: %v", err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(uLogin.Password), []byte(loginRequest.Password))
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, api.ErrBadRequest(fmt.Errorf("compare password: %v", err), ""))
		return
	}

	token, err := token.ClaimToken(uLogin.Username)
	if err != nil {
		u.Log.Printf("ERROR : %+v", err)
		api.ResponseError(w, fmt.Errorf("claim token: %v", err))
		return
	}

	var res response.TokenResponse
	res.Token = token

	api.ResponseOK(w, res, http.StatusOK)
}
