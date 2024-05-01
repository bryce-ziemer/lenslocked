package main

import (
	"fmt"

	"bryce-ziemer/github.com/lenslocked/models"
)

// varies depending on mail service
const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 2525
	username = "d2d80aa2bea15a"
	password = "6d22ab77885635"
)

func main() {

	email := models.Email{
		From:      "test@lenslocked.com",
		To:        "bryce.ziemer@gmail.com",
		Subject:   "This is a test email",
		Plaintext: "This is the body of the email",
		HTML:      `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`,
	}

	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})

	err := es.Send(email)
	if err != nil {
		panic(err)
	}

	fmt.Println("Msg Sent")

}
