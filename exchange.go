package main

import (
	"errors"
	"fmt"
	"io"
	email "net/mail"
)

type Channel interface {
	Hostname() string
	Reply(code int, msg string)
}

type WriterChannel struct {
	w    io.Writer
	host string
}

func (c *WriterChannel) Hostname() string {
	return c.host
}

func (c *WriterChannel) Reply(code int, msg string) {
	reply := fmt.Sprintf("%d %s\r\n", code, msg)
	c.w.Write([]byte(reply))
}

type Exchange struct {
	io.Reader
	Channel

	domain string
	from   *email.Address
	to     []*email.Address
	body   io.Reader

	mailer Mailer
	parser *email.AddressParser
}

func NewExchange(m Mailer, r io.Reader, c Channel) *Exchange {
	return &Exchange{
		Channel: c,
		Reader:  r,
		mailer:  m,
		parser:  &email.AddressParser{},
	}
}

func (ex *Exchange) Domain(domain string) error {
	if !isDomainName(domain) && !isIp(domain) {
		return errors.New("invalid domain")
	}

	ex.domain = domain
	return nil
}

func (ex *Exchange) From(from string) error {
	addr, err := ex.parser.Parse(from)
	if err != nil {
		return err
	}

	ex.from = addr
	return nil
}

func (ex *Exchange) To(to string) error {
	addr, err := ex.parser.Parse(to)
	if err != nil {
		return err
	}

	if ex.to == nil {
		ex.to = []*email.Address{addr}
	} else {
		ex.to = append(ex.to, addr)
	}

	return nil
}

func (ex *Exchange) Body(r io.Reader) {
	ex.body = r
}

func (ex *Exchange) Done() {
	ex.mailer.Send(Mail{
		From: ex.from,
		To:   ex.to,
		Body: ex.body,
	})

	ex.Reset()
}

func (ex *Exchange) Reset() {
	ex.from = nil
	ex.to = nil
}
