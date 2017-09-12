package main

import (
	"errors"
	"fmt"
	"io"
	"net/mail"
	"strings"
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
	lines := strings.Split(msg, "\n")

	for i, line := range lines {
		last := i == len(lines)-1

		space := " "
		if !last {
			space = "-"
		}

		out := fmt.Sprintf("%d%s%s\r\n", code, space, line)
		c.w.Write([]byte(out))
	}
}

type Exchange struct {
	io.ReadCloser
	Channel

	domain string
	from   *mail.Address
	to     []*mail.Address
	body   io.Reader

	mailer Mailer
	parser *mail.AddressParser
}

func NewExchange(m Mailer, r io.ReadCloser, c Channel) *Exchange {
	return &Exchange{
		Channel:    c,
		ReadCloser: r,
		mailer:     m,
		parser:     &mail.AddressParser{},
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
		ex.to = []*mail.Address{addr}
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
