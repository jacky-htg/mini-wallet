package controllers

import (
	"database/sql"
	"mini-wallet/libraries/api"
	"mini-wallet/libraries/database"
	"net/http"
)

// Checks : struct for set Checks Dependency Injection
type Checks struct {
	Db *sql.DB
}

// Health : http handler for login
func (u Checks) Health(w http.ResponseWriter, r *http.Request) {
	var health struct {
		Status string `json:"healthy"`
	}

	// Check if the database is ready.
	if err := database.StatusCheck(r.Context(), u.Db); err != nil {

		// If the database is not ready we will tell the client and use a 500 status.
		health.Status = "db not ready"
		api.ResponseError(w, err)
		return
	}

	health.Status = "ok"
	api.ResponseOK(w, health, http.StatusOK)
}
