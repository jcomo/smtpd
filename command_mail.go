package main

type mailCommand struct{}

func (c *mailCommand) Next() []string {
	return []string{CommandRcpt, CommandRset}
}

func (c *mailCommand) Process(line string, ex *Exchange) (*reply, bool) {
	address, got := getSuffix(line, "MAIL FROM:")
	if !got {
		r := newReply(ReplySyntaxError, "Syntax: MAIL FROM: <address>")
		return r, false
	}

	ex.Reset()
	err := ex.From(address)
	if err != nil {
		r := newReply(ReplyInvalidMailboxSyntax, err.Error())
		return r, false
	}

	return ok(), true
}
