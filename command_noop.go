package main

type noopCommand struct{}

func (c *noopCommand) Next() []string {
	return nil
}

func (c *noopCommand) Process(line string, ex *Exchange) (bool, error) {
	ex.Reply(ReplyOK, "noop")
	return true, nil
}
