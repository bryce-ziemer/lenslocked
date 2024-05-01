package main

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

// varies depending on mail service
const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 2525
	username = "d2d80aa2bea15a"
	password = "6d22ab77885635"
)

func main() {

	from := "test@lenslocked.com"
	to := "bryce.ziemer@gmail.com" //"jon@calhoun.io"
	subject := "This is a test email"
	plaintext := "This is the body of the email"
	html := `<h1?Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`

	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plaintext)
	msg.AddAlternative("text/html", html)

	//msg.WriteTo(os.Stdout)

	dialer := mail.NewDialer(host, port, username, password)
	err := dialer.DialAndSend(msg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Msg Sent")

}
