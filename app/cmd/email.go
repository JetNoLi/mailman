package main

import (
	"ewails/app/server/config"
	"ewails/app/server/services"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/mail"
)

// var EmailServerAddr = "imap.gmail.com:993"
var EmailServerAddr = "imap.secureserver.net:993"

func ProcessMail(c *imapclient.Client, seqNum uint32) (imap.FetchItemBodySection, error) {
	// Send a FETCH command to fetch the message body
	seqSet := imap.SeqSetNum(seqNum)
	bodySection := &imap.FetchItemBodySection{}
	fetchOptions := &imap.FetchOptions{
		BodySection: []*imap.FetchItemBodySection{bodySection},
	}
	fetchCmd := c.Fetch(seqSet, fetchOptions)
	defer fetchCmd.Close()

	msg := fetchCmd.Next()
	if msg == nil {
		return imap.FetchItemBodySection{}, fmt.Errorf("FETCH command did not return any message")
	}

	// Find the body section in the response
	var bodySectionData imapclient.FetchItemDataBodySection
	ok := false
	for {
		item := msg.Next()
		if item == nil {
			break
		}
		bodySectionData, ok = item.(imapclient.FetchItemDataBodySection)
		if ok {
			break
		}
	}
	if !ok {
		return imap.FetchItemBodySection{}, fmt.Errorf("FETCH command did not return body section")
	}

	// Read the message via the go-message library
	mr, err := mail.CreateReader(bodySectionData.Literal)
	if err != nil {
		return imap.FetchItemBodySection{}, fmt.Errorf("failed to create mail reader: %v", err)
	}

	// Print a few header fields
	h := mr.Header
	if date, err := h.Date(); err != nil {
		log.Printf("failed to parse Date header field: %v", err)
	} else {
		log.Printf("Date: %v", date)
	}
	if to, err := h.AddressList("To"); err != nil {
		log.Printf("failed to parse To header field: %v", err)
	} else {
		log.Printf("To: %v", to)
	}
	if subject, err := h.Text("Subject"); err != nil {
		log.Printf("failed to parse Subject header field: %v", err)
	} else {
		log.Printf("Subject: %v", subject)
	}

	// Process the message's parts
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return imap.FetchItemBodySection{}, fmt.Errorf("failed to read message part: %v", err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, _ := io.ReadAll(p.Body)
			log.Printf("Inline text: %v", string(b))
		case *mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			log.Printf("Attachment: %v", filename)
		}
	}

	if err := fetchCmd.Close(); err != nil {
		return imap.FetchItemBodySection{}, fmt.Errorf("FETCH command failed: %v", err)
	}

	return *bodySection, nil
}

func GetEmail() {
	config.ReadEnv()

	var Email = os.Getenv("TEST_EMAIL")
	var Pswd = os.Getenv("TEST_EMAIL_PSWD")

	// Install Imap
	client, err := imapclient.DialTLS(EmailServerAddr, &imapclient.Options{})

	if err != nil {
		log.Fatal("error creating client", err.Error())
	}

	cap, err := client.Capability().Wait()

	if err != nil {
		log.Fatal("error getting capabalities", err.Error())
	}

	fmt.Println(cap.AuthMechanisms(), cap.QuotaResourceTypes())

	login := client.Login(Email, Pswd)

	if err = login.Wait(); err != nil {
		log.Fatal("error logging in", err.Error())
	}

	//TODO: Defer close
	mailboxes, err := client.List("", "%", &imap.ListOptions{
		ReturnStatus: &imap.StatusOptions{
			NumMessages: true,
			NumUnseen:   true,
		},
	}).Collect()

	if err != nil {
		log.Fatalf("error getting mailboxes %s", err.Error())
	}

	var mBox *imap.ListData

	for _, mailbox := range mailboxes {
		fmt.Println(mailbox.Mailbox)

		if mailbox.Mailbox == "INBOX" {
			mBox = mailbox
		}
	}

	if mBox == nil {
		log.Fatalf("no mailbox selected")
	}

	log.Printf("Mailbox %q contains %v messages (%v unseen)", mBox.Mailbox, *mBox.Status.NumMessages, *mBox.Status.NumUnseen)

	selectData, err := client.Select(mBox.Mailbox, &imap.SelectOptions{
		ReadOnly: true,
	}).Wait()

	if err != nil {
		log.Fatalf("error selecting inbox %s", err.Error())
	}

	fmt.Println(selectData)

	data, err := client.UIDSearch(&imap.SearchCriteria{
		Body: []string{"katlego"},
	}, nil).Wait()

	if err != nil {
		log.Fatalf("UID SEARCH command failed: %v", err)
	}

	emailIds := data.AllUIDs()
	emailId := emailIds[0]

	fmt.Println("email id yeah", emailId)

	email, err := ProcessMail(client, 95)

	if err != nil {
		log.Fatalf("error processing mail: %s", err.Error())

	}

	fmt.Println(email)

}

func Read() {
	config.ReadEnv()

	es, err := services.NewEmailService(EmailServerAddr)

	if err != nil {
		log.Fatal(err.Error())
	}

	var Email = os.Getenv("TEST_EMAIL")
	var Pswd = os.Getenv("TEST_EMAIL_PSWD")

	err = es.Login(Email, Pswd)

	if err != nil {
		log.Fatal(err.Error())
	}

	mailboxes, err := es.ListMailboxes()

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("mailboxes, ", mailboxes)

	err = es.ConnectToMailbox("INBOX", true)

	if err != nil {
		log.Fatal(err.Error())
	}

	msgs, err := es.ListMail(5)

	if err != nil {
		log.Fatal("error listing messages", err)
	}

	fmt.Println(msgs)
}

func main() {
	// GetEmail()
	Read()
}
