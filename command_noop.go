package main

var noop Command = Command{
	Name:      CommandNoop,
	Stateless: true,
	Run: func(line string, ex *Exchange) bool {
		ex.Reply(ReplyOK, "noop")
		return true
	},
}
