package main

type heloCommand struct{}

func (c *heloCommand) Next() []string {
	return []string{CommandMail}
}

func (c *heloCommand) Process(line string, ex *Exchange) (*reply, bool) {
	domain, got := getSuffix(line, "HELO ")
	if !got {
		r := newReply(ReplySyntaxError, "Syntax: HELO <domain>")
		return r, false
	}

	err := ex.Domain(domain)
	if err != nil {
		r := newReply(ReplySyntaxError, err.Error())
		return r, false
	}

	r := newReply(ReplyOK, "HELO "+ex.Hostname())
	return r, true
}
