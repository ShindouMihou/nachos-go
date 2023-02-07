package nachos

import "github.com/nats-io/nats.go"

type Context struct {
	Message *nats.Msg
	Store   map[string]any
}

func (ctx *Context) Get(key string) any {
	if res, exs := ctx.Store[key]; !exs {
		return nil
	} else {
		return res
	}
}

func (ctx *Context) GetOrDefault(key string, def any) any {
	if res, exs := ctx.Store[key]; !exs {
		return def
	} else {
		return res
	}
}
