package smtpd

import (
	"fmt"
	"io"
	"net/mail"
	"os"
)

type Mailer interface {
	Send(msg *mail.Message) error
}

type DebugMailer struct{}

func (m *DebugMailer) Send(msg *mail.Message) error {
	fmt.Println("From: " + msg.Header.Get("From"))
	fmt.Println("To: " + msg.Header.Get("To"))
	fmt.Println()

	io.Copy(os.Stdout, msg.Body)
	return nil
}
