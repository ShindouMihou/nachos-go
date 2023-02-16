### 〰️ nachos-go

A minimal and simple routing abstraction over Nats.io that enables a more HTTP-like routing system with special features 
such as "before-action" and "end-action" which all contribute to making the NATS.io simpler than it already is.

Nachos-go is intentionally single-instance, therefore, the library cannot support projects that have more than one application. This 
is to make the overall developer experience of Nachos-go simpler.

#### 🌃 Creating Routes

To create routes in nachos-go, you can simply a new `nachos.Route` and subscribe it:
```go
package main
func main() {
	route := nachos.Route{
		Path: "quests",
		Children: []nachos.Route{
			{
				Path: "accept",
				Action: func (ctx *nachos.Context) {
					nachos.Reply(ctx.Message, "Some message here")
				},
			},
			{
				Path: "delete",
				Action: func (ctx *nachos.Context) {
					nachos.Reply(ctx.Message, struct { Message string `json:"message"` }{Message:"Acknowledged"})
				},
			},
		},
	}
	nachos.Subscribe(route)
}
```

You can add before and end actions by simply adding them to the route declaration, if the route has children nodes then 
all the before and end actions will be inherited by the children nodes. An example of a before action would be:
```go
package main
func main() {
	route := nachos.Route{
		Path: "quest.delete",
		BeforeAction: []BeforeActionEvent{
			func(ctx *nachos.Context) bool {
				if ctx.Message.Header.Get("Magic-Word") !== "Please" {
					nachos.Reply(ctx.Message, "Invalid magic word, please try better!")
					return false
				}
				return true
			},
		},
		Action: func (ctx *nachos.Context) {
			nachos.Reply(ctx.Message, "Yay, you sent the right magic word!")
		},
	}
	nachos.Subscribe(route)
	// ... additional stuff here
}
```

#### 🚃 Queue Groups

By default, nachos-go doesn't subscribe by queue, but rather a global one, which is generally not the desired behavior of 
many applications. In nachos-go, we have two methods to enable queue grouping and in the following hierarchy:
1. Per-route Queue Group
2. Default Queue Group

To add a queue group to a single route, you can fill the Queue Group property, such in the following example:
```go
package main

import "github.com/ShindouMihou/nachos-go/nachos"

func main() {
	route := nachos.Route{
		// ... additional properties here
		QueueGroup: &nachos.QueueGroup{
			Enabled: true,
			Name: "nachos:go",
        },
	}
	nachos.Subscribe(route)
	// ... additional stuff here
}
```

To add a default queue group, you can override the `nachos.DefaultQueueGroup` property such in the example:

```go
package main

import "github.com/ShindouMihou/nachos-go/nachos"

func main() {
	nachos.DefaultQueueGroup = &nachos.QueueGroup{
		Enabled: true,
		Name: "nachos:go",
    }
}
```

#### 🧰 Troubleshoot

> **Unable to subscribe more routes after doing `nachos.Attach()`**
>
>nachos-go doesn't support adding additional routes after running `nachos.Attach()` and that is because nachos holds
almost no reference to the connection. Currently, there is no point in supporting this functionality, therefore,
it will remain unsupported for the time-being.