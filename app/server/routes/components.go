package routes

import (
	"ewails/app/server/handlers"
	"ewails/wrappers/rtr"
)

func CompRouter() *rtr.Router {
	r := rtr.New("/")

	r.Get("/home/main-menu", handlers.GetMainMenuOptions)
	emailRtr := rtr.New("/email/")

	emailRtr.Get("/mailbox/menu", handlers.GetMailboxMenu)
	emailRtr.Get("/account/menu", handlers.GetAccountMenu)

	r.Use("/", emailRtr)

	return r
}
