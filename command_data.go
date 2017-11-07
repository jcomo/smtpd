package smtpd

import (
	"bytes"
	"io"
	"strings"
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
	return CommandsMail
}

func (c *dataCommand) Process(line string, ex *Exchange) (*reply, bool) {
	if !c.started {
		c.started = true
		r := newReply(ReplyDataStart,
			"start mail input; end with <CRLF>.<CRLF>")
		return r, false
	}

	if strings.HasPrefix(line, ".") {
		// We remove the first dot since clients will send a preceding
		// dot to avoid the EOM sequence
		line = line[1:]

		if len(line) == 0 {
			ex.Body(c.reader())
			err := ex.Done()

			if err != nil {
				c.buf.Reset()
				r := newReply(ReplySyntaxError, err.Error())
				return r, false
			}

			return ok(), true
		}
	}

	c.buf.WriteString(line)
	c.buf.WriteRune('\n')
	return nil, false
}

func (c *dataCommand) reader() io.Reader {
	// Transparently remove the last byte since it is an extraneous newline
	stream := c.buf.Bytes()
	stream = stream[:len(stream)-1]
	return bytes.NewReader(stream)
}
