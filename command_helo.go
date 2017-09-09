package main

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
