package main

type rsetCommand struct{}

func (c *rsetCommand) Next() []string {
	return []string{CommandMail}
}

func (c *rsetCommand) Process(line string, ex *Exchange) (*reply, bool) {
	ex.Reset()
	return ok(), true
}
