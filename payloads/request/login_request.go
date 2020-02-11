package request

import (
	"errors"
	"mini-wallet/libraries/api"
	"net/http"
)

// LoginRequest : format json request for login
type LoginRequest struct {
	Username string
	Password string
}

// Transform form to LoginRequest. Include Form Validate.
func (u *LoginRequest) Transform(r *http.Request) error {
	if len(r.FormValue("username")) <= 0 {
		return api.ErrBadRequest(errors.New("name is required"), "name is required")
	}

	if len(r.FormValue("password")) <= 0 {
		return api.ErrBadRequest(errors.New("password is required"), "password is required")
	}

	u.Username = r.FormValue("username")
	u.Password = r.FormValue("password")

	return nil
}
