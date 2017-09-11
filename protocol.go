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

func defaultCommands() map[string]commandFactory {
	return map[string]commandFactory{
		CommandData: factory(newDataCommand),
		CommandEhlo: nil,
		CommandExpn: nil,
		CommandHelo: instanceFactory(&heloCommand{}),
		CommandHelp: nil,
		CommandMail: instanceFactory(&mailCommand{}),
		CommandNoop: instanceFactory(&noopCommand{}),
		CommandRcpt: instanceFactory(&rcptCommand{}),
		CommandRset: instanceFactory(&rsetCommand{}),
		CommandSaml: nil,
		CommandSoml: nil,
		CommandVrfy: nil,
	}
}
