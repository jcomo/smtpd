package main

type mailCommand struct{}

func (c *mailCommand) Next() []string {
	return []string{CommandRcpt, CommandRset}
}

func (c *mailCommand) Process(line string, ex *Exchange) (bool, error) {
	address, got := getSuffix(line, "MAIL FROM:")
	if !got {
		ex.Reply(ReplySyntaxError, "Syntax: MAIL FROM: <address>")
		return false, nil
	}

	ex.Reset()
	err := ex.From(address)
	if err != nil {
		ex.Reply(ReplyInvalidMailboxSyntax, err.Error())
		return false, nil
	}

	ex.Reply(ReplyOK, "OK")
	return true, nil
}
