package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	html "html/template"
	"log"
	"net/smtp"
	"strconv"
	"text/template"
)

// Service interface to users service
type Service interface {
	SendText(recipients []string, subject, text string) error
	SendInvite(orgName, code string, recipients []string) error
	// TODO: make templates for matches etc
}

type mailService struct {
	config ClientConfig
	addr   string
	auth   smtp.Auth
}

// NewMailService will return a struct that implements the mailService interface
func NewMailService(config ClientConfig) *mailService {
	return &mailService{
		config: config,
		addr:   config.Host + ":" + strconv.Itoa(config.Port),
		auth:   smtp.PlainAuth("", config.Username, config.Password, config.Host),
	}
}

// SendText sends a basic plaintext email
func (s *mailService) SendText(recipients []string, subject, text string) error {
	var body bytes.Buffer

	templateData := struct {
		Body string
	}{
		Body: text,
	}

	t, err := template.New("textTemplate").Parse(`{{.Body}}`)
	if err != nil {
		log.Print("error trying to parse mail template")
		return err
	}

	err = t.Execute(&body, templateData)
	if err != nil {
		log.Print("error trying to execute mail template")
		return err
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	subject = fmt.Sprintf("Subject: %s\n", subject)
	msg := []byte(subject + mime + "\n" + body.String())

	err = smtp.SendMail(s.addr, s.auth, s.config.Username, recipients, msg)
	if err != nil {
		log.Printf("error sending mail: %s\n", err)
		return err
	}

	return nil
}

// SendInvite sends an invite email
func (s *mailService) SendInvite(orgName, code string, recipients []string) error {
	codeBase64 := base64.StdEncoding.EncodeToString([]byte(code))
	link := fmt.Sprintf("http://localhost:3000/invite/%s", codeBase64)

	subject := "Invite"
	templatePath := "./mail/templates/invite.html"

	templateData := struct {
		Link         string
		Organization string
	}{
		Link:         link,
		Organization: orgName,
	}

	err := s.sendHTML(recipients, subject, templatePath, templateData)
	if err != nil {
		log.Printf("error sending html email template: %s\n", err)
		return err
	}

	return nil
}

func (s *mailService) sendHTML(recipients []string, subject, path string, data interface{}) error {
	var body bytes.Buffer

	template, err := html.ParseFiles(path)
	if err != nil {
		log.Printf("error creating email template: %s\n", err)
		return err
	}

	err = template.Execute(&body, data)
	if err != nil {
		log.Print("error trying to execute mail template")
		return err
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject = fmt.Sprintf("Subject: %s\n", subject)
	msg := []byte(subject + mime + "\n" + body.String())

	err = smtp.SendMail(s.addr, s.auth, s.config.Username, recipients, msg)
	if err != nil {
		log.Printf("error sending mail: %s\n", err)
		return err
	}

	return nil
}
