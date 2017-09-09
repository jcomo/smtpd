package main

import (
	"bufio"
	"bytes"
)

var data Command = Command{
	Name: CommandData,
	Next: []string{
		CommandMail,
	},
	Run: func(line string, ex *Exchange) bool {
		ex.Reply(ReplyDataStart, "start mail input; end with <CRLF>.<CRLF>")

		var buf bytes.Buffer
		scanner := bufio.NewScanner(ex)

		for {
			scanner.Scan()
			bs := scanner.Bytes()
			if len(bs) > 0 && bs[0] == EOD {
				break
			}

			buf.Write(bs)
			buf.WriteByte(NEWLINE)
		}

		ex.Body(bytes.NewReader(buf.Bytes()))
		ex.Done()

		ex.Reply(ReplyOK, "OK")
		return true
	},
}
