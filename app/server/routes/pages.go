package routes

import (
	"ewails/app/pages"
	"ewails/app/server/handlers"
	"ewails/wrappers/rtr"
)

func PagesRouter() *rtr.Router {
	r := rtr.New("/")

	r.Get("/landing", handlers.TemplHandler(pages.Index()))

	r.Get("/login", handlers.TemplHandler(pages.Login()))

	r.Get("/home", handlers.TemplHandler(pages.Home()))

	return r
}
