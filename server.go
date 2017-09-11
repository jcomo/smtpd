package main

import (
	"io"
)

type Server struct {
	Hostname string
	Mailer   Mailer
	IOLoop   IOLoop

	commands map[string]commandFactory
}

func NewServer(host string, loop IOLoop) *Server {
	return &Server{
		Hostname: host,
		Mailer:   &DebugMailer{},
		IOLoop:   loop,
		commands: defaultCommands(),
	}
}

func (s *Server) Run() {
	err := s.IOLoop.Run(s.accept)
	if err != nil {
		panic(err)
	}
}

func (s *Server) accept(r io.ReadCloser, w io.Writer) {
	channel := &WriterChannel{w: w, host: s.Hostname}
	exchange := NewExchange(s.Mailer, r, channel)
	loop := &protocolLoop{
		ex:       exchange,
		commands: s.commands,
		next:     []string{CommandHelo},
		waiting:  true,
	}

	loop.Run()
}
