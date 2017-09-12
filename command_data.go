package main

import (
	"bytes"
)

func newDataCommand() command {
	return &dataCommand{
		started: false,
	}
}

type dataCommand struct {
	started bool
	buf     bytes.Buffer
}

func (c *dataCommand) Next() []string {
	return []string{CommandMail}
}

func (c *dataCommand) Process(line string, ex *Exchange) (*reply, bool) {
	if !c.started {
		c.started = true
		r := newReply(ReplyDataStart,
			"start mail input; end with <CRLF>.<CRLF>")
		return r, false
	}

	if line == "." {
		stream := c.buf.Bytes()
		ex.Body(bytes.NewReader(stream))
		ex.Done()

		return ok(), true
	}

	c.buf.WriteString(line)
	c.buf.WriteRune('\n')
	return nil, false
}
