package services

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/mail"
)

type EmailService struct {
	client *imapclient.Client
}

// TODO: Create opts funcs later
func NewEmailService(imapServer string) (EmailService, error) {
	client, err := imapclient.DialTLS(imapServer, &imapclient.Options{})

	if err != nil {
		return EmailService{}, err
	}

	return EmailService{
		client: client,
	}, nil
}

func ParseEmail(data *imapclient.FetchMessageData) (Email, error) {
	if data == nil {
		return Email{}, fmt.Errorf("no message found")
	}

	// Find the body section in the response
	var bodySectionData imapclient.FetchItemDataBodySection
	ok := false
	for {
		item := data.Next()
		if item == nil {
			break
		}
		bodySectionData, ok = item.(imapclient.FetchItemDataBodySection)
		if ok {
			break
		}
	}

	if !ok {
		return Email{}, fmt.Errorf("no message body found")
	}

	mr, err := mail.CreateReader(bodySectionData.Literal)

	if err != nil {
		return Email{}, err
	}

	e := Email{}

	h := mr.Header
	if date, err := h.Date(); err != nil {
		log.Printf("failed to parse Date header field: %v", err)
	} else {
		e.Date = date.GoString()
	}
	if to, err := h.AddressList("To"); err != nil {
		log.Printf("failed to parse To header field: %v", err)
	} else {
		for _, recipient := range to {
			e.Recipients = append(e.Recipients, recipient.Address)
		}
	}
	if subject, err := h.Text("Subject"); err != nil {
		log.Printf("failed to parse Subject header field: %v", err)
	} else {
		e.Subject = subject
	}

	fmt.Println("parsing body", e)

	messageData := ""

	attachments := []EmailAttachment{}

	// // Process the message's parts
	for {
		p, err := mr.NextPart()
		if errors.Is(err, io.EOF) {
			fmt.Println("jumping due to body")
			break
		} else if err != nil {
			fmt.Printf("failed to read message part: %v", err)
			break
		}

		fmt.Println("running")

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, err := io.ReadAll(p.Body)

			if err != nil {
				return e, err
			}

			log.Printf("Inline text: %v", string(b))
			messageData += string(b)
		case *mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			log.Printf("Attachment: %v", filename)
			data, err := io.ReadAll(p.Body)

			if err != nil {
				return e, err
			}

			attachments = append(attachments, EmailAttachment{
				Name: filename,
				Data: data,
			})
		}
	}

	e.Attatchments = attachments
	e.Body = messageData

	return e, nil
}

func (es *EmailService) Login(email string, password string) error {
	return es.client.Login(email, password).Wait()
}

type Mailbox struct {
	Id      string
	Name    string
	Total   uint32
	Unread  uint32
	NextUid uint32
}

// TODO: Create and Figure Out Options
func (es *EmailService) ListMailboxes() ([]Mailbox, error) {
	listCmd := es.client.List("", "%", &imap.ListOptions{
		ReturnStatus: &imap.StatusOptions{
			NumMessages: true,
			NumUnseen:   true,
		},
	})

	mailBoxes := []Mailbox{}

	for mailbox := listCmd.Next(); mailbox != nil; mailbox = listCmd.Next() {
		if mailbox.Status == nil {
			break
		}

		mailBoxes = append(mailBoxes, Mailbox{
			Id:     mailbox.Mailbox,
			Name:   mailbox.Mailbox,
			Total:  *mailbox.Status.NumMessages,
			Unread: *mailbox.Status.NumUnseen,
		})
	}

	if err := listCmd.Close(); err != nil {
		return nil, err
	}

	return mailBoxes, nil
}

// TODO: Figure out options and optsfuncs
func (es *EmailService) ConnectToMailbox(mailboxName string, readOnly bool) error {
	_, err := es.client.Select(mailboxName, &imap.SelectOptions{
		ReadOnly: readOnly,
	}).Wait()

	return err
}

// TODO: Rather Use Client Status Command
func (es *EmailService) GetCurrentMailBox() (Mailbox, error) {
	mailbox := es.client.Mailbox()

	if mailbox == nil {
		return Mailbox{}, fmt.Errorf("no mailbox selected")
	}

	status, err := es.client.Status(mailbox.Name, &imap.StatusOptions{
		NumMessages: true,
		NumUnseen:   true,
		UIDNext:     true,
	}).Wait()

	if err != nil {
		return Mailbox{}, err
	}

	return Mailbox{
		Id:      status.Mailbox,
		Total:   *status.NumMessages,
		Unread:  *status.NumUnseen,
		Name:    status.Mailbox,
		NextUid: uint32(status.UIDNext),
	}, nil
}

type EmailAttachment struct {
	Name string
	Data []byte
}

type Email struct {
	Body         string // html body string
	Date         string
	Subject      string
	Recipients   []string
	From         string
	Attatchments []EmailAttachment
}

func (es *EmailService) ListMail(prev uint32) ([]Email, error) {
	mailBox, err := es.GetCurrentMailBox()

	if err != nil {
		return nil, err
	}

	currentUid := mailBox.NextUid - 2

	if prev > currentUid {
		prev = currentUid - 1
	}

	start := currentUid
	end := currentUid - prev

	uids := make([]uint32, start-end)

	for index := range start - end {
		id := start - index
		fmt.Println(id, index)
		uids[index] = id
	}

	return es.FetchManyById(uids)
}

func (es *EmailService) FetchById(uid uint32) (Email, error) {
	bodySection := &imap.FetchItemBodySection{}
	fetchOptions := &imap.FetchOptions{
		Flags:       true,
		Envelope:    true,
		BodySection: []*imap.FetchItemBodySection{bodySection},
	}

	fetchCmd := es.client.Fetch(imap.SeqSetNum(uid), fetchOptions)

	msg := fetchCmd.Next()

	if err := fetchCmd.Close(); err != nil {
		return Email{}, err
	}

	return ParseEmail(msg)
}

func (es *EmailService) FetchManyById(uids []uint32) ([]Email, error) {
	bodySection := &imap.FetchItemBodySection{}
	fetchOptions := &imap.FetchOptions{
		Flags:       true,
		Envelope:    true,
		BodySection: []*imap.FetchItemBodySection{bodySection},
	}

	fmt.Println("uids", uids)
	fetchCmd := es.client.Fetch(imap.SeqSetNum(uids...), fetchOptions)

	emails := []Email{}
	errs := []error{}

	for {
		fmt.Println("looping")
		msg := fetchCmd.Next()

		if msg == nil {
			break
		}

		email, err := ParseEmail(msg)

		if err != nil {
			errs = append(errs, err)
			continue
		}

		emails = append(emails, email)
	}

	fmt.Println("fetch errs", errs)

	if err := fetchCmd.Close(); err != nil {
		return nil, err
	}

	return emails, nil
}
