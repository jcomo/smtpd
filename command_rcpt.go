package main

type rcptCommand struct{}

func (c *rcptCommand) Next() []string {
	return []string{CommandData, CommandRcpt, CommandRset}
}

func (c *rcptCommand) Process(line string, ex *Exchange) (bool, error) {
	address, got := getSuffix(line, "RCPT TO:")
	if !got {
		ex.Reply(ReplySyntaxError, "Syntax: RCPT TO: <address>")
		return false, nil
	}

	err := ex.To(address)
	if err != nil {
		ex.Reply(ReplyInvalidMailboxSyntax, err.Error())
		return false, nil
	}

	ex.Reply(ReplyOK, "OK")
	return true, nil
}
