package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
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

		err := ex.Domain(domain)
		if err != nil {
			ex.Reply(ReplySyntaxError, err.Error())
			return false
		}

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
	if !isDomainName(domain) {
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

type WriterChannel struct {
	w io.Writer
}

func (c *WriterChannel) Reply(code int, msg string) {
	reply := fmt.Sprintf("%d %s\n", code, msg)
	c.w.Write([]byte(reply))
}

func contains(a []string, s string) bool {
	for _, e := range a {
		if e == s {
			return true
		}
	}

	return false
}

// Lifted from the net pkg in the go std library. This function is private.
// Instead of importing an entire DNS library for one function, we just
// use it directly.
func isDomainName(s string) bool {
	// See RFC 1035, RFC 3696.
	// Presentation format has dots before every label except the first, and the
	// terminal empty label is optional here because we assume fully-qualified
	// (absolute) input. We must therefore reserve space for the first and last
	// labels' length octets in wire format, where they are necessary and the
	// maximum total length is 255.
	// So our _effective_ maximum is 253, but 254 is not rejected if the last
	// character is a dot.
	l := len(s)
	if l == 0 || l > 254 || l == 254 && s[l-1] != '.' {
		return false
	}

	last := byte('.')
	ok := false // Ok once we've seen a letter.
	partlen := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		default:
			return false
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_':
			ok = true
			partlen++
		case '0' <= c && c <= '9':
			// fine
			partlen++
		case c == '-':
			// Byte before dash cannot be dot.
			if last == '.' {
				return false
			}
			partlen++
		case c == '.':
			// Byte before dot cannot be dot, dash.
			if last == '.' || last == '-' {
				return false
			}
			if partlen > 63 || partlen == 0 {
				return false
			}
			partlen = 0
		}
		last = c
	}
	if last == '-' || partlen > 63 {
		return false
	}
	return ok
}

func run(ex *Exchange) {
	next := []string{CommandHelo}
	scanner := bufio.NewScanner(ex)

	ex.Reply(ReplyServiceReady, "Simple Mail Transfer Service Ready")

	for {
		got := scanner.Scan()
		if !got {
			err := scanner.Err()
			if err != nil {
				log.Println(err)
			}

			break
		}

		line := scanner.Text()
		name := safeSubstring(line, 4)

		if name == CommandQuit {
			ex.Reply(ReplyServiceClosing, "Service closing transmission channel")
			break
		}

		if contains(unimplemented, name) {
			ex.Reply(ReplyNotImplemented, "not implemented")
			continue
		}

		cmd, ok := commands[name]
		if !ok {
			ex.Reply(ReplyUnknown, "unrecognized command")
			continue
		}

		if cmd.Stateless {
			cmd.Run(line, ex)
			continue
		}

		if !contains(next, name) {
			ex.Reply(ReplyBadSequence, "bad command sequence")
			continue
		}

		if cmd.Run(line, ex) {
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

func handleConn(conn net.Conn) {
	mailer := &DebugMailer{}
	channel := &WriterChannel{w: conn}
	exchange := NewExchange(mailer, conn, channel)

	run(exchange)
}

var smtpHostVar string
var smtpPortVar int

func main() {
	flag.StringVar(&smtpHostVar, "smtp-host", "127.0.0.1",
		"The host on which to run the SMTP server")
	flag.IntVar(&smtpPortVar, "smtp-port", 8025,
		"The port on which to run the SMTP server")

	flag.Parse()
	smtpAddr := fmt.Sprintf("%s:%d", smtpHostVar, smtpPortVar)

	l, err := net.Listen("tcp", smtpAddr)
	if err != nil {
		panic(err)
	}

	log.Printf("SMTP server listening on %s\n", smtpAddr)

	for {
		// TODO: graceful shutdown with draining
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		go handleConn(conn)
	}
}
