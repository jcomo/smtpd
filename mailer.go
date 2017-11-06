package smtpd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/mail"
	"os"
)

type Mail struct {
	From *mail.Address
	To   []*mail.Address
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

type addrPayload struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address"`
}

type mailPayload struct {
	From addrPayload   `json:"from"`
	To   []addrPayload `json:"to"`
	Body string        `json:"body"`
}

func createPayload(mail Mail) *mailPayload {
	from := addrPayload{
		Name:    mail.From.Name,
		Address: mail.From.Address,
	}

	to := []addrPayload{}
	for _, addr := range mail.To {
		to = append(to, addrPayload{
			Name:    addr.Name,
			Address: addr.Address,
		})
	}

	buf := &bytes.Buffer{}
	io.Copy(buf, mail.Body)

	return &mailPayload{
		From: from,
		To:   to,
		Body: buf.String(),
	}
}

type HTTPMailer struct {
	url    string
	client *http.Client
}

func NewHTTPMailer(url string) *HTTPMailer {
	return &HTTPMailer{
		url:    url,
		client: http.DefaultClient,
	}
}

func (m *HTTPMailer) Send(mail Mail) error {
	data, err := json.Marshal(createPayload(mail))
	if err != nil {
		return nil
	}

	buf := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, m.url, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	log.Println("Sending mail to " + m.url)
	_, err = m.client.Do(req)
	return err
}
