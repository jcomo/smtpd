package main

import (
	"flag"
	"log"
	"os"
)

var debugVar bool
var smtpHostVar string
var smtpPortVar int

func main() {
	flag.BoolVar(&debugVar, "debug", false,
		"Run in debug mode (do NOT use in production)")
	flag.StringVar(&smtpHostVar, "smtp-host", "localhost",
		"The host on which to run the SMTP server")
	flag.IntVar(&smtpPortVar, "smtp-port", 8025,
		"The port on which to run the SMTP server")

	flag.Parse()
	host, err := os.Hostname()
	if err != nil {
		log.Println("No hostname available. Using defined host instead")
		host = smtpHostVar
	}

	var loop IOLoop
	if debugVar {
		loop = &ConsoleIO{}
	} else {
		loop = NewSocketIO(smtpHostVar, smtpPortVar)
	}

	srv := NewServer(host, loop)
	srv.Run()
}
