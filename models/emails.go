package models

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@lenslocked.com"
)

type Email struct {
	From      string
	To        string
	Subject   string
	Plaintext string
	HTML      string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailService(config SMTPConfig) *EmailService {
	es := EmailService{
		dialer: mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}

	return &es
}

// Good use of service because heavily relies on 3rd party library??
type EmailService struct {
	// DefaultSender is used as the default sender when one isn't provided for an
	// email. This is also used in functions where the email is predetermined,
	// like the forgotten password email.
	DefaultSender string

	//unexported fields
	dialer *mail.Dialer
}

func (es EmailService) Send(email Email) error {
	// instructor said we wrap the new message function and create a message using Email struct
	// so other devs do not need to know about the mail dependency...
	msg := mail.NewMessage()
	msg.SetHeader("To", email.To)
	es.setFrom(msg, email)
	msg.SetHeader("Subject", email.Subject)

	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.Plaintext != "" && email.HTML != "":
		msg.SetHeader("Subject", email.Subject)
		msg.AddAlternative("text/html", email.HTML)
	default:
		return fmt.Errorf("send: Empty Body")
	}

	err := es.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("send: %w", err)

	}
	return nil
}

func (es *EmailService) ForgotPassword(to string, resetURL string) error {
	email := Email{
		Subject:   "Reset your password",
		To:        to,
		Plaintext: "To reset your password please visit the following link: " + resetURL,
		HTML: `<p>To reset your password, please visit the following link: 
		<a href="` + resetURL + `">` + resetURL + `</a><p>`,
	}

	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("ForgotPassword email: %w", err)
	}
	return nil
}

func (es *EmailService) setFrom(msg *mail.Message, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case es.DefaultSender != "":
		from = es.DefaultSender
	default:
		from = DefaultSender
	}
	msg.SetHeader("From", from)
}
