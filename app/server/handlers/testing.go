package handlers

import (
	"ewails/app/pages"
	"ewails/app/server/services"
	"net/http"
	"os"
)

func TestingLoginSkip(w http.ResponseWriter, r *http.Request) {
	es, err := services.NewEmailService("imap.secureserver.net:993")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	testEmail := os.Getenv("TEST_EMAIL")

	err = es.Login(testEmail, os.Getenv("TEST_EMAIL_PSWD"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	es.SetAccount(&services.EmailAccount{
		Addr:       testEmail,
		ImapServer: "imap.secureserver.net:993",
	})

	err = SetEmailService(r, &es)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = RenderTempl(&w, r, pages.Home())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
