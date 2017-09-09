package main

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

const (
	EOD     = 0x2e // '.'
	NEWLINE = 0xa  // \n
)

var commands map[string]Command
var unimplemented []string

type RunFunc func(string, *Exchange) bool

type Command struct {
	Name      string
	Next      []string
	Stateless bool
	Run       RunFunc
}

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
