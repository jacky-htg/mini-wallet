package api

import (
	"context"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// App struct for application
type App struct {
	mux *httprouter.Router
	log *log.Logger
	mw  []Middleware
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

// Handler custome httprouter
type Handler func(http.ResponseWriter, *http.Request)

// Ctx for index context
type Ctx string

// Handle httprouter
func (a *App) Handle(method, url string, h Handler) {
	h = wrapMiddleware(a.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), Ctx("ps"), ps)
		ctx = context.WithValue(ctx, Ctx("url"), url)
		h(w, r.WithContext(ctx))
	}

	a.mux.Handle(method, url, fn)
}

// HandleCors and OPTIONS response
func (a *App) HandleCors() {
	a.mux.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Add("Access-Control-Allow-Origin", "*")
			header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS, PUT")
			header.Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, Token")
			header.Add("Content-Type", "application/json; charset=utf-8")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})
}

// NewApp for create new app
func NewApp(log *log.Logger, mw ...Middleware) *App {
	return &App{
		mux: httprouter.New(),
		log: log,
		mw:  mw,
	}
}
