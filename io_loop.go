package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type AcceptFunc func(io.ReadCloser, io.Writer)

type IOLoop interface {
	Run(AcceptFunc) error
}

type ConsoleIO struct{}

func (cio *ConsoleIO) Run(accept AcceptFunc) error {
	log.Println("!!! smtpd is running in debug mode. " +
		"This should NEVER be enabled in production")

	accept(os.Stdin, os.Stdout)
	return nil
}

type SocketIO struct {
	addr string
}

func NewSocketIO(host string, port int) *SocketIO {
	addr := fmt.Sprintf("%s:%d", host, port)
	return &SocketIO{addr: addr}
}

func (sio *SocketIO) Run(accept AcceptFunc) error {
	l, err := net.Listen("tcp", sio.addr)
	if err != nil {
		return err
	}

	log.Printf("SMTP server listening on %s\n", sio.addr)

	for {
		// TODO: graceful shutdown with draining
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go accept(conn, conn)
	}

	return nil
}
