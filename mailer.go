package main

import (
	"fmt"
	"io"
	email "net/mail"
	"os"
)

type Mail struct {
	From *email.Address
	To   []*email.Address
	Body io.Reader
}

type Mailer interface {
	Send(mail Mail) error
}

type DebugMailer struct {
}

func (m *DebugMailer) Send(mail Mail) error {
	fmt.Println("From: " + mail.From.String())
	for _, addr := range mail.To {
		fmt.Println("To: " + addr.String())
	}

	fmt.Println()
	io.Copy(os.Stdout, mail.Body)
	return nil
}
