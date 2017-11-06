package smtpd

type rsetCommand struct{}

func (c *rsetCommand) Next() []string {
	return CommandsMail
}

func (c *rsetCommand) Process(line string, ex *Exchange) (*reply, bool) {
	ex.Reset()
	return ok(), true
}
