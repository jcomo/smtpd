package main

type rsetCommand struct{}

func (c *rsetCommand) Next() []string {
	return []string{CommandMail}
}

func (c *rsetCommand) Process(line string, ex *Exchange) (bool, error) {
	ex.Reset()
	ex.Reply(ReplyOK, "reset")
	return true, nil
}
