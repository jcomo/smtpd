package main

var mail Command = Command{
	Name: CommandMail,
	Next: []string{
		CommandRcpt,
		CommandRset,
	},
	Run: func(line string, ex *Exchange) bool {
		address, got := getSuffix(line, "MAIL FROM:")
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
