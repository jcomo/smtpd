package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

func handleConn(host string, conn net.Conn) {
	mailer := &DebugMailer{}
	channel := &WriterChannel{w: conn, host: host}
	handle(NewExchange(mailer, conn, channel))
}

func handle(ex *Exchange) {
	next := []string{CommandHelo}
	scanner := bufio.NewScanner(ex)

	ex.Reply(ReplyServiceReady, "Simple Mail Transfer Service Ready")

	for {
		got := scanner.Scan()
		if !got {
			err := scanner.Err()
			if err != nil {
				log.Println(err)
			}

			break
		}

		line := scanner.Text()
		name := strings.ToUpper(safeSubstring(line, 4))

		if name == CommandQuit {
			ex.Reply(ReplyServiceClosing, "Service closing transmission channel")
			break
		}

		if contains(unimplemented, name) {
			ex.Reply(ReplyNotImplemented, "not implemented")
			continue
		}

		cmd, ok := commands[name]
		if !ok {
			ex.Reply(ReplyUnknown, "unrecognized command")
			continue
		}

		if cmd.Stateless {
			cmd.Run(line, ex)
			continue
		}

		if !contains(next, name) {
			ex.Reply(ReplyBadSequence, "bad command sequence")
			continue
		}

		if cmd.Run(line, ex) {
			next = cmd.Next
		}
	}
}
