package main

type reply struct {
	Code    int
	Message string
}

func newReply(code int, msg string) *reply {
	return &reply{Code: code, Message: msg}
}

func ok() *reply {
	return newReply(ReplyOK, "OK")
}

type command interface {
	Next() []string
	Process(line string, ex *Exchange) (*reply, bool)
}

type commandFactory interface {
	New() command
}

type factoryFunc struct {
	factory func() command
}

func (f *factoryFunc) New() command {
	return f.factory()
}

func instanceFactory(c command) *factoryFunc {
	return &factoryFunc{
		factory: func() command {
			return c
		},
	}
}

func factory(f func() command) *factoryFunc {
	return &factoryFunc{factory: f}
}
