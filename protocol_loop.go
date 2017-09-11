package main

import (
	"bufio"
	"log"
	"strings"
)

type protocolLoop struct {
	ex       *Exchange
	commands map[string]commandFactory
	command  command
	next     []string
	waiting  bool
}

func (l *protocolLoop) Run() {
	defer l.ex.Close()

	scanner := bufio.NewScanner(l.ex)
	l.ex.Reply(ReplyServiceReady, "Simple Mail Transfer Service Ready")

	for {
		got := scanner.Scan()
		line := scanner.Text()

		if !got {
			err := scanner.Err()
			if err != nil {
				log.Println(err)
			}

			break
		}

		if !l.waiting {
			l.processCommand(line)
			continue
		}

		name := l.readCommandName(line)
		if name == CommandQuit {
			l.ex.Reply(ReplyServiceClosing, "Service closing transmission channel")
			break
		}

		factory, ok := l.commands[name]
		if !ok {
			l.ex.Reply(ReplyUnknown, "unrecognized command")
			continue
		}

		if factory == nil {
			l.ex.Reply(ReplyNotImplemented, "not implemented")
			continue
		}

		command := factory.New()
		if !l.canRun(name, command) {
			l.ex.Reply(ReplyBadSequence, "bad command sequence")
			continue
		}

		l.command = command
		l.processCommand(line)
	}
}

func (l *protocolLoop) processCommand(line string) {
	done, err := l.command.Process(line, l.ex)
	if err != nil {
		log.Println(err)
		l.waiting = true
		return
	}

	l.waiting = done
	if done {
		next := l.command.Next()
		if next != nil {
			l.next = next
		}
	}
}

func (l *protocolLoop) canRun(name string, c command) bool {
	return l.isStateless(c) || contains(l.next, name)
}

func (l *protocolLoop) isStateless(c command) bool {
	return c.Next() == nil
}

func (l *protocolLoop) readCommandName(line string) string {
	return strings.ToUpper(safeSubstring(line, 4))
}
