package main

import (
	"bufio"
	"fmt"
	"io"
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
		scanner := wordScanner(line)
		scanner.Scan()
		got := scanner.Scan()
		if !got {
			ex.Reply(ReplySyntaxError, "Syntax: HELO <domain>")
			return false
		}

		domain := scanner.Text()
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

		ex.From(address)
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

		ex.To(address)
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

func getSuffix(s, prefix string) (string, bool) {
	if !strings.HasPrefix(s, prefix) {
		return "", false
	} else {
		return s[len(prefix):], true
	}
}

func wordScanner(s string) *bufio.Scanner {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)
	return scanner
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
	Channel
	domain string
	from   string
	to     []string
}

func NewExchange(c Channel) *Exchange {
	return &Exchange{Channel: c}
}

func (ex *Exchange) Domain(domain string) error {
	// TODO: check form
	ex.domain = domain
	return nil
}

func (ex *Exchange) From(from string) error {
	// TODO: check form
	ex.from = from
	return nil
}

func (ex *Exchange) To(to string) error {
	// TODO: check form
	if ex.to == nil {
		ex.to = []string{to}
	} else {
		ex.to = append(ex.to, to)
	}

	return nil
}

func (ex *Exchange) Reset() {
	ex.from = ""
	ex.to = nil
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
	exchange := NewExchange(c)

	c.Reply(ReplyServiceReady, "Simple Mail Transfer Service Ready")

	for {
		fmt.Println(exchange)
		scanner.Scan() // TODO: error handling, eof etc
		text := scanner.Text()
		if len(text) < 4 {
			c.Reply(ReplyUnknown, "no command")
			continue
		}

		name := strings.ToUpper(text[:4])
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

func main() {
	run(os.Stdin, &StdChannel{})
}
