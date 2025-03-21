package handlers

import (
	"ewails/app/components/assets"
	"ewails/app/components/auth"
	"net/http"
	"os"
)

const (
	GoogleOAuthProviderUrl = "https://accounts.google.com/o/oauth2/v2/auth"
)

// func LoginWithEmailAndPassword(w http.ResponseWriter, r *http.Request) {
// 	LoginButton := auth.LoginWithEmailAndPassword()
// 	err := RenderTempl(&w, r, LoginButton)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

func GoogleOAuthInitiate(w http.ResponseWriter, r *http.Request) {
	LoginButton := auth.LoginWithGoogle(os.Getenv("CLIENT_ID"), os.Getenv("GOOGLE_REDIRECT_URI"))
	err := RenderTempl(&w, r, LoginButton)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func MailMainImage(w http.ResponseWriter, r *http.Request) {
	mailMan := assets.MailMan()
	err := RenderTempl(&w, r, mailMan)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GoogleOAuthCallback(w http.ResponseWriter, r *http.Request) {

}
