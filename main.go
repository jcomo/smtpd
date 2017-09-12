package main

import (
	"flag"
	"fmt"
)

var debugVar bool
var hookUrlVar string
var smtpHostVar string
var smtpPortVar int

func usage() {
	fmt.Printf("usage: smtpd [options]\n\n")
	fmt.Printf("A server with a compatible implementation of RFC 851.\n\n")
	fmt.Printf("Flags:\n")

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("  --%-16s %s\n", f.Name, f.Usage)
	})
}

func main() {
	flag.BoolVar(&debugVar, "debug", false,
		"Run in debug mode (do NOT use in production)")
	flag.StringVar(&hookUrlVar, "hook-url", "",
		"If specified, sent mail will be delivered here via HTTP POST")
	flag.StringVar(&smtpHostVar, "smtp-host", "localhost",
		"The network interface to bind to")
	flag.IntVar(&smtpPortVar, "smtp-port", 8025,
		"Change the port (default: 8025)")

	flag.Usage = usage
	flag.Parse()

	var loop IOLoop
	var mailer Mailer

	loop = &ConsoleIO{}
	if !debugVar {
		loop = NewSocketIO(smtpHostVar, smtpPortVar)
	}

	mailer = &DebugMailer{}
	if hookUrlVar != "" {
		mailer = NewHTTPMailer(hookUrlVar)
	}

	srv := NewServer()
	srv.IOLoop = loop
	srv.Mailer = mailer
	srv.Run()
}
