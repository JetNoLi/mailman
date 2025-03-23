package server

import (
	"ewails/app/server/config"
	"ewails/app/server/routes"
	"ewails/wrappers/rtr"
	"fmt"
	"net/http"
)

func Create(middleware ...rtr.MiddlewareFn) *rtr.Router {
	config.ReadEnv()

	r := rtr.New("/")

	r.UseMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getting request")
	})

	r.UseMiddleware(middleware...)

	r.Use("/auth/", routes.AuthRouter())
	r.Use("/pages/", routes.PagesRouter())
	r.Use("/components/", routes.CompRouter())

	return r
}
