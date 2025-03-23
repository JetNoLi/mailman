package handlers

import (
	"ewails/app/components/assets"
	"ewails/app/components/auth"
	"ewails/app/pages"
	"ewails/app/server/services"
	"fmt"
	"net/http"
	"os"
)

const (
	GoogleOAuthProviderUrl = "https://accounts.google.com/o/oauth2/v2/auth"
)

func LoginWithEmailAndPassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "error parsing form "+err.Error(), http.StatusBadRequest)
		return
	}

	//TODO: Validate
	email := r.PostForm.Get("email")
	pswd := r.PostForm.Get("password")
	imapAddr := "imap.secureserver.net:993" //TODO: Move to form

	es, err := services.NewEmailService(imapAddr)

	if err != nil {
		http.Error(w, fmt.Sprintf("error creating email service %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = es.Login(email, pswd)

	if err != nil {
		http.Error(w, fmt.Sprintf("error logging in %s", err.Error()), http.StatusInternalServerError)
		return
	}

	es.SetAccount(&services.EmailAccount{Addr: email, ImapServer: imapAddr})

	err = SetEmailService(r, &es)

	if err != nil {
		http.Error(w, fmt.Sprintf("error logging in %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Render Page
	Home := pages.Home()
	err = RenderTempl(&w, r, Home)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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
