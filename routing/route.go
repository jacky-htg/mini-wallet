package routing

import (
	"database/sql"
	"log"
	"mini-wallet/controllers"
	"mini-wallet/libraries/api"
	"mini-wallet/middlewares"
	"net/http"
)

func mid(db *sql.DB, log *log.Logger) []api.Middleware {
	prefix := "/api/v1"

	var mw []api.Middleware
	mw = append(mw, middlewares.Auths(db, log, []string{
		"POST " + prefix + "/login",
		"POST " + prefix + "/init",
		"GET " + prefix + "/health",
	}))

	return mw
}

// API handler routing
func API(db *sql.DB, log *log.Logger) http.Handler {
	prefix := "/api/v1"

	app := api.NewApp(log, mid(db, log)...)
	app.HandleCors()

	// Health Routing
	{
		check := controllers.Checks{Db: db}
		app.Handle(http.MethodGet, prefix+"/health", check.Health)
	}

	// Auth Routing
	{
		auth := controllers.Auths{Db: db, Log: log}
		app.Handle(http.MethodPost, prefix+"/login", auth.Login)
	}

	// Wallet Routing
	{
		wallet := controllers.Wallets{Db: db, Log: log}
		app.Handle(http.MethodGet, prefix+"/wallet", wallet.View)
		app.Handle(http.MethodPost, prefix+"/wallet", wallet.Enabled)
		app.Handle(http.MethodPatch, prefix+"/wallet", wallet.Disabled)
		app.Handle(http.MethodPost, prefix+"/wallet/deposits", wallet.Deposit)
		app.Handle(http.MethodPost, prefix+"/wallet/withdrawals", wallet.Withdrawal)
		app.Handle(http.MethodPost, prefix+"/init", wallet.Init)
	}

	return app
}
