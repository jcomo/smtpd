package main

type heloCommand struct{}

func (c *heloCommand) Next() []string {
	return []string{CommandMail}
}

func (c *heloCommand) Process(line string, ex *Exchange) (bool, error) {
	domain, got := getSuffix(line, "HELO ")
	if !got {
		ex.Reply(ReplySyntaxError, "Syntax: HELO <domain>")
		return false, nil
	}

	err := ex.Domain(domain)
	if err != nil {
		ex.Reply(ReplySyntaxError, err.Error())
		return false, nil
	}

	ex.Reply(ReplyOK, "HELO "+ex.Hostname())
	return true, nil
}
