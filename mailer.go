package smtpd

import (
	"fmt"
	"io"
	"net/mail"
	"os"
)

type Mailer interface {
	Send(mail *mail.Message) error
}

type DebugMailer struct{}

func (m *DebugMailer) Send(mail *mail.Message) error {
	fmt.Println("From: " + mail.Header.Get("From"))
	fmt.Println("To: " + mail.Header.Get("To"))
	fmt.Println()

	io.Copy(os.Stdout, mail.Body)
	return nil
}
