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

func (c *dataCommand) Process(line string, ex *Exchange) (bool, error) {
	if !c.started {
		ex.Reply(ReplyDataStart, "start mail input; end with <CRLF>.<CRLF>")
		c.started = true
		return false, nil
	}

	if line == "." {
		stream := c.buf.Bytes()
		ex.Body(bytes.NewReader(stream))
		ex.Done()

		ex.Reply(ReplyOK, "OK")
		return true, nil
	}

	c.buf.WriteString(line)
	c.buf.WriteRune('\n')
	return false, nil
}
