package routes

import (
	"ewails/app/pages"
	"ewails/app/server/handlers"
	"ewails/wrappers/rtr"
	"net/http"
)

func PagesRouter() *rtr.Router {
	r := rtr.New("/")

	r.Get("/landing", func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderTempl(&w, r, pages.Index())
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderTempl(&w, r, pages.Login())
	})

	return r
}
