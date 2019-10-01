package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	t "html/template"
	"log"
	"net/smtp"
	"strconv"
	"text/template"
)

type SmtpTemplateData struct {
	From    string
	To      string
	Subject string
	Body    string
}

// Service interface to users service
type Service interface {
	SendText(msg string, recipients []string) error
	SendInvite(orgName, code string, recipient string) error
}

type mailService struct {
	config ClientConfig
}

// NewMailService will return a struct that implements the mailService interface
func NewMailService(config ClientConfig) *mailService {
	return &mailService{config: config}
}

const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}
`

func (s *mailService) SendText(msg string, recipients []string) error {
	var err error
	var doc bytes.Buffer

	context := &SmtpTemplateData{
		"lunchmoreapp@gmail.com",
		"liam.tj.pimlott@gmail.com",
		"Hello Subject",
		msg,
	}

	t := template.New("emailTemplate")
	t, err = t.Parse(emailTemplate)
	if err != nil {
		log.Print("error trying to parse mail template")
		return err
	}
	err = t.Execute(&doc, context)
	if err != nil {
		log.Print("error trying to execute mail template")
		return err
	}

	from := "lunchmoreapp@gmail.com"

	// log.Printf("creds: user:%s pass:%s\n", s.config.Username, s.config.Password)
	// hostname is used by PlainAuth to validate the TLS certificate.
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	err = smtp.SendMail(s.config.Host+":"+strconv.Itoa(s.config.Port), auth, from, recipients, doc.Bytes())
	if err != nil {
		log.Printf("error sending mail: %s\n", err)
		return err
	}
	return nil
}

type inviteData struct {
	From         string
	To           string
	Subject      string
	Link         string
	Organization string
}

const inviteTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

<a href="{{.Link}}">Click here</a> to join {{.Organization}}!
`

func (s *mailService) SendInvite(orgName, code string, recipient string) error {
	var doc bytes.Buffer
	from := "lunchmoreapp@gmail.com"

	// make link
	codeBase64 := base64.StdEncoding.EncodeToString([]byte(code))
	link := fmt.Sprintf("https://localhost:8080/invite/%s", codeBase64)

	data := &inviteData{
		from,
		recipient,
		fmt.Sprintf("Invite to %s on Lunchmore", orgName),
		link,
		orgName,
	}

	template, err := t.New("invite email").Parse(inviteTemplate)
	if err != nil {
		log.Printf("error creating email template: %s\n", err)
		return err
	}

	err = template.Execute(&doc, data)
	if err != nil {
		log.Print("error trying to execute mail template")
		return err
	}

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	to := []string{recipient}
	err = smtp.SendMail(s.config.Host+":"+strconv.Itoa(s.config.Port), auth, from, to, doc.Bytes())
	if err != nil {
		log.Printf("error sending mail: %s\n", err)
		return err
	}

	return nil
}
