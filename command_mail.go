package main

import "fmt"

type mailCommand struct {
	name string
}

func newMailCommand(name string) *mailCommand {
	return &mailCommand{name: name}
}

func (c *mailCommand) Next() []string {
	return []string{CommandRcpt, CommandRset}
}

func (c *mailCommand) Process(line string, ex *Exchange) (*reply, bool) {
	prefix := fmt.Sprintf("%s FROM:", c.name)
	address, got := getSuffix(line, prefix)
	if !got {
		r := newReply(ReplySyntaxError,
			fmt.Sprintf("Syntax: %s <address>", prefix))
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
