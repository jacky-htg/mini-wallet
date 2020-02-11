package middlewares

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"mini-wallet/libraries/api"
	"mini-wallet/libraries/array"
	"mini-wallet/models"
)

// Auths middleware
func Auths(db *sql.DB, log *log.Logger, allow []string) api.Middleware {
	fn := func(before api.Handler) api.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			var isAuth bool
			var str array.ArrString
			var customer models.Customer
			isAuth = true

			ctx := r.Context()
			curRoute := ctx.Value(api.Ctx("url")).(string)
			curRoute = strings.ToUpper(r.Method) + " " + curRoute

			inArray, _ := str.InArray(curRoute, allow)
			if !inArray {
				var err error
				url := r.URL.String()
				controller := strings.Split(url, "/")[1]

				err = customer.IsAuth(
					ctx,
					db,
					r.Header.Get("Authorization"),
					controller,
					curRoute,
				)

				if err == sql.ErrNoRows {
					isAuth = false
					log.Printf("ERROR : %+v", err)
					api.ResponseError(w, api.ErrForbidden(errors.New("Forbidden"), ""))
					return
				}

				if err != nil {
					isAuth = false
					log.Printf("ERROR : %+v", err)
					api.ResponseError(w, err)
					return
				}
			}

			if !isAuth {
				log.Print("ERROR : Forbidden")
				api.ResponseError(w, api.ErrForbidden(errors.New("Forbidden"), ""))
				return
			}

			ctx = context.WithValue(ctx, api.Ctx("auth"), customer)
			before(w, r.WithContext(ctx))
		}

		return h
	}

	return fn
}
