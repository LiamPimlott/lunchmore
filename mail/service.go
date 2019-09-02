package mail

import (
	"bytes"
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

const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}
`

// Service interface to users service
type Service interface {
	SendText(msg string, recipients []string) error
}

type mailService struct {
	config ClientConfig
}

// NewMailService will return a struct that implements the mailService interface
func NewMailService(config ClientConfig) *mailService {
	return &mailService{config: config}
}

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
	}
	err = t.Execute(&doc, context)
	if err != nil {
		log.Print("error trying to execute mail template")
	}

	from := "lunchmoreapp@gmail.com"

	log.Printf("creds: user:%s pass:%s\n", s.config.Username, s.config.Password)

	// hostname is used by PlainAuth to validate the TLS certificate.
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	err = smtp.SendMail(s.config.Host+":"+strconv.Itoa(s.config.Port), auth, from, recipients, doc.Bytes())
	if err != nil {
		log.Printf("error sending mail: %s\n", err)
		return err
	}
	return nil
}
