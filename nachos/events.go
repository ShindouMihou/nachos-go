package nachos

type BeforeActionEvent = func(ctx *Context) bool
type ActionEvent = func(ctx *Context)
type EndAction = func(ctx *Context)

// EndActionEvent defines the general structure of a middleware such as how it should operate.
// You can configure the EndActionEvent to ignore the result of a middleware by setting Passthrough to true
// which would make the EndActionEvent execute regardless of whether the ActionEvent executed.
type EndActionEvent struct {
	Passthrough bool
	Action      EndAction
}
