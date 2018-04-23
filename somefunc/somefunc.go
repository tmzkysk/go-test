package somefunc

type Caller interface {
	call(val int) int
}

type Client struct {
	FuncCaller Caller
}

type ExampleCaller struct{}

func (c *Client) Run(val int) int {
	return c.FuncCaller.call(val)
}

func (f *ExampleCaller) call(val int) int {
	return val
}
