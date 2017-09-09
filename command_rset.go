package main

var rset Command = Command{
	Name: CommandRset,
	Next: []string{
		CommandMail,
	},
	Run: func(line string, ex *Exchange) bool {
		ex.Reset()
		ex.Reply(ReplyOK, "reset")
		return true
	},
}
