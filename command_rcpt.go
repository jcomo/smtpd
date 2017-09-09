package main

var rcpt Command = Command{
	Name: CommandRcpt,
	Next: []string{
		CommandData,
		CommandRcpt,
		CommandRset,
	},
	Run: func(line string, ex *Exchange) bool {
		address, got := getSuffix(line, "RCPT TO:")
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
