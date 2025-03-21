package server

import (
	"ewails/app/server/config"
	"ewails/app/server/routes"
	"ewails/wrappers/rtr"
	"fmt"
	"net/http"
)

func Create() *rtr.Router {
	config.ReadEnv()

	r := rtr.New("/")

	r.UseMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getting request")
	})

	r.Use("/auth/", routes.AuthRouter())
	r.Use("/pages/", routes.PagesRouter())

	return r
}
