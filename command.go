package main

type command interface {
	Next() []string
	Process(line string, ex *Exchange) (bool, error)
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
