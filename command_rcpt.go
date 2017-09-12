package main

type rcptCommand struct{}

func (c *rcptCommand) Next() []string {
	return []string{CommandData, CommandRcpt, CommandRset}
}

func (c *rcptCommand) Process(line string, ex *Exchange) (*reply, bool) {
	address, got := getSuffix(line, "RCPT TO:")
	if !got {
		r := newReply(ReplySyntaxError, "Syntax: RCPT TO: <address>")
		return r, false
	}

	err := ex.To(address)
	if err != nil {
		r := newReply(ReplyInvalidMailboxSyntax, err.Error())
		return r, false
	}

	return ok(), true
}
