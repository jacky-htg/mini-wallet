package models

import (
	"context"
	"database/sql"
	"errors"
	"mini-wallet/libraries/api"
	"mini-wallet/libraries/token"
	"strings"
)

// Customer : struct of Customer
type Customer struct {
	ID       string
	Name     string
	Username string
	Password string
	Email    string
}

// GetByUsername : get user by username
func (u *Customer) GetByUsername(ctx context.Context, db *sql.DB) error {
	const q string = `SELECT customer_id, name, username, password, email FROM customers`
	err := db.QueryRowContext(ctx, q+" WHERE username = $1", u.Username).Scan(&u.ID, &u.Name, &u.Username, &u.Password, &u.Email)

	if err == sql.ErrNoRows {
		err = api.ErrNotFound(err, "")
	}

	return err
}

// IsAuth for check customer authorization
func (u *Customer) IsAuth(ctx context.Context, db *sql.DB, tokenparam interface{}, controller string, route string) error {
	if tokenparam == nil {
		return api.ErrBadRequest(errors.New("Bad request for token"), "")
	}

	tokenString := tokenparam.(string)
	if tokens := strings.Split(tokenString, " "); len(tokens) >= 2 {
		tokenString = tokens[1]
	}

	isValid, username := token.ValidateToken(tokenString)
	if !isValid {
		return api.ErrBadRequest(errors.New("Bad request for invalid token"), "")
	}

	u.Username = username

	return u.GetByUsername(ctx, db)
}
