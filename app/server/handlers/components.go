package handlers

import (
	"ewails/app/components/email"
	"ewails/app/components/home"
	"fmt"
	"net/http"
)

func GetMailboxMenu(w http.ResponseWriter, r *http.Request) {
	es, err := GetEmailService(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mailboxes, err := es.ListMailboxes()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = RenderTempl(&w, r, email.MailBoxMenu("INBOX", mailboxes))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetAccountMenu(w http.ResponseWriter, r *http.Request) {
	es, err := GetEmailService(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	account := es.GetAccountInfo()

	if account.Addr == "" {
		http.Error(w, "no account details specified", http.StatusInternalServerError)
		return
	}

	err = RenderTempl(&w, r, email.AccountInfoMenu(account.Addr, account.ImapServer))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetMainMenuOptions(w http.ResponseWriter, r *http.Request) {
	es, err := GetEmailService(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mailboxes, err := es.ListMailboxes()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var total, totalUnread uint32 = 0, 0

	for _, mailbox := range mailboxes {
		total += mailbox.Total
		totalUnread += mailbox.Unread
	}

	totalTxt := fmt.Sprintf("%d Total Emails", total)
	totalUnreadTxt := fmt.Sprintf("%d Total Unread Emails", totalUnread)

	options := []home.MenuOption{
		{Title: "Quick Search", SubText: totalTxt, Asset: "src/assets/images/search.svg", OnClickHxGet: "/components/home/menu/search/"},
		{Title: "Quick Sort", SubText: totalUnreadTxt, Asset: "src/assets/images/Folders.svg"},
		{Title: "Subscription Manager", SubText: "", Asset: "src/assets/images/marketing.svg"},
	}

	err = RenderTempl(&w, r, home.MainMenuOptions(options))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetSearchMenu(w http.ResponseWriter, r *http.Request) {
	err := RenderTempl(&w, r, home.SearchMenu())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
