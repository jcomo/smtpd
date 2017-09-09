package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var smtpHostVar string
var smtpPortVar int

func main() {
	flag.StringVar(&smtpHostVar, "smtp-host", "localhost",
		"The host on which to run the SMTP server")
	flag.IntVar(&smtpPortVar, "smtp-port", 8025,
		"The port on which to run the SMTP server")

	flag.Parse()
	smtpAddr := fmt.Sprintf("%s:%d", smtpHostVar, smtpPortVar)

	l, err := net.Listen("tcp", smtpAddr)
	if err != nil {
		panic(err)
	}

	log.Printf("SMTP server listening on %s\n", smtpAddr)

	for {
		// TODO: graceful shutdown with draining
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("new connection from " + conn.RemoteAddr().String())
		go handleConn(conn)
	}
}
