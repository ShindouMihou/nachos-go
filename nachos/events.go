package nachos

type BeforeActionEvent = func(ctx *Context) bool
type ActionEvent = func(ctx *Context)
type EndActionEvent = func(ctx *Context)
