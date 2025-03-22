package routes

import (
	"ewails/app/server/handlers"
	"ewails/wrappers/rtr"
)

func AuthRouter() *rtr.Router {

	authRouter := rtr.New("/")

	authRouter.Post("/login/", handlers.LoginWithEmailAndPassword)

	authRouter.Post("/mailman/", handlers.MailMainImage)

	//OAuth
	oauthRouter := rtr.New("/oauth")

	oauthRouter.Get("/google/initiate", handlers.GoogleOAuthInitiate)
	oauthRouter.Post("/google/callback", handlers.GoogleOAuthCallback)

	return authRouter
}
