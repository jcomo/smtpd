package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	email "net/mail"
	"os"
	"strings"
)

const (
	ReplyServiceReady         = 220
	ReplyServiceClosing       = 221
	ReplyOK                   = 250
	ReplyDataStart            = 354
	ReplyUnknown              = 500
	ReplySyntaxError          = 501
	ReplyNotImplemented       = 502
	ReplyBadSequence          = 503
	ReplyInvalidMailboxSyntax = 553
)

const (
	CommandData = "DATA"
	CommandEhlo = "EHLO"
	CommandExpn = "EXPN"
	CommandHelo = "HELO"
	CommandHelp = "HELP"
	CommandMail = "MAIL"
	CommandNoop = "NOOP"
	CommandQuit = "QUIT"
	CommandRcpt = "RCPT"
	CommandRset = "RSET"
	CommandSaml = "SAML"
	CommandSoml = "SOML"
	CommandVrfy = "VRFY"
)

type RunFunc func(string, *Exchange) bool

type Command struct {
	Name      string
	Next      []string
	Stateless bool
	Run       RunFunc
}

var helo Command = Command{
	Name: CommandHelo,
	Next: []string{
		CommandMail,
	},
	Run: func(line string, ex *Exchange) bool {
		domain, got := getSuffix(line, "HELO ")
		if !got {
			ex.Reply(ReplySyntaxError, "Syntax: HELO <domain>")
			return false
		}

		ex.Domain(domain)
		ex.Reply(ReplyOK, "HELO "+domain)
		return true
	},
}

var mail Command = Command{
	Name: CommandMail,
	Next: []string{
		CommandRcpt,
		CommandRset,
	},
	Run: func(line string, ex *Exchange) bool {
		address, got := getSuffix(line, "MAIL FROM: ")
		if !got {
			ex.Reply(ReplySyntaxError, "Syntax: MAIL FROM: <address>")
			return false
		}

		ex.Reset()
		err := ex.From(address)
		if err != nil {
			ex.Reply(ReplyInvalidMailboxSyntax, err.Error())
			return false
		}

		ex.Reply(ReplyOK, "OK")
		return true
	},
}

var rcpt Command = Command{
	Name: CommandRcpt,
	Next: []string{
		CommandData,
		CommandRcpt,
		CommandRset,
	},
	Run: func(line string, ex *Exchange) bool {
		address, got := getSuffix(line, "RCPT TO: ")
		if !got {
			ex.Reply(ReplySyntaxError, "Syntax: RCPT TO: <address>")
			return false
		}

		err := ex.To(address)
		if err != nil {
			ex.Reply(ReplyInvalidMailboxSyntax, err.Error())
			return false
		}

		ex.Reply(ReplyOK, "OK")
		return true
	},
}

var rset Command = Command{
	Name: CommandRset,
	Next: []string{
		CommandMail,
	},
	Run: func(line string, ex *Exchange) bool {
		ex.Reset()
		ex.Reply(ReplyOK, "reset")
		return true
	},
}

var noop Command = Command{
	Name:      CommandNoop,
	Stateless: true,
	Run: func(line string, ex *Exchange) bool {
		ex.Reply(ReplyOK, "noop")
		return true
	},
}

const (
	EOD     = 0x2e // '.'
	NEWLINE = 0xa  // \n
)

var data Command = Command{
	Name: CommandData,
	Next: []string{
		CommandMail,
	},
	Run: func(line string, ex *Exchange) bool {
		ex.Reply(ReplyDataStart, "start mail input; end with <CRLF>.<CRLF>")

		var buf bytes.Buffer
		scanner := bufio.NewScanner(ex)

		for {
			scanner.Scan()
			bs := scanner.Bytes()
			if len(bs) > 0 && bs[0] == EOD {
				break
			}

			// TODO: error handling / checking
			buf.Write(bs)
			buf.WriteByte(NEWLINE)
		}

		ex.Body(bytes.NewReader(buf.Bytes()))
		ex.Done()

		ex.Reply(ReplyOK, "OK")
		return true
	},
}

func getSuffix(s, prefix string) (string, bool) {
	if !strings.HasPrefix(s, prefix) {
		return "", false
	} else {
		return s[len(prefix):], true
	}
}

var commands map[string]Command
var unimplemented []string

func init() {
	commands = make(map[string]Command)
	commands[helo.Name] = helo
	commands[mail.Name] = mail
	commands[rcpt.Name] = rcpt
	commands[rset.Name] = rset
	commands[noop.Name] = noop
	commands[data.Name] = data

	unimplemented = []string{
		CommandHelp,
		CommandEhlo,
		CommandExpn,
		CommandVrfy,
		CommandSoml,
		CommandSaml,
	}
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
	// TODO: check form
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

type Channel interface {
	Reply(code int, msg string)
}

type StdChannel struct{}

func (c *StdChannel) Reply(code int, msg string) {
	fmt.Printf("%d %s\n", code, msg)
}

func contains(a []string, s string) bool {
	for _, e := range a {
		if e == s {
			return true
		}
	}

	return false
}

func run(r io.Reader, c Channel) {
	next := []string{CommandHelo}
	scanner := bufio.NewScanner(r)
	mailer := &DebugMailer{}
	exchange := NewExchange(mailer, r, c)

	c.Reply(ReplyServiceReady, "Simple Mail Transfer Service Ready")

	for {
		scanner.Scan() // TODO: error handling, eof etc
		text := scanner.Text()
		name := safeSubstring(text, 4)

		if name == CommandQuit {
			c.Reply(ReplyServiceClosing, "Service closing transmission channel")
			break
		}

		if contains(unimplemented, name) {
			c.Reply(ReplyNotImplemented, "not implemented")
			continue
		}

		cmd, ok := commands[name]
		if !ok {
			c.Reply(ReplyUnknown, "unrecognized command")
			continue
		}

		if cmd.Stateless {
			cmd.Run(text, exchange)
			continue
		}

		if !contains(next, name) {
			c.Reply(ReplyBadSequence, "bad command sequence")
			continue
		}

		if cmd.Run(text, exchange) {
			next = cmd.Next
		}
	}
}

func safeSubstring(s string, n int) string {
	if len(s) < n {
		return s
	} else {
		return s[:n]
	}
}

func main() {
	run(os.Stdin, &StdChannel{})
}
