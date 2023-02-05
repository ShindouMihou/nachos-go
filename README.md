### nachos-go

A small and simple abstraction on top of NATS.io that enables a more HTTP-like routing system with special features 
such as "route-condition", "before-action" and "end-action" which all contribute to making the NATS.io simpler than it 
already is.

Nachos-go is intentionally single-instance, therefore, the library cannot support projects that have more than one application. This 
is to make the overall developer experience of Nachos-go simpler.

#### Creating Routes

To create routes in nachos-go, you can simply a new `nachos.Route` and subscribe it:
```go
route := nachos.Route{
	Path: "quests",
	Children: []nachos.Route{
		{
			Path: "accept",
			Action: func (message *nats.Msg) {
				nachos.Reply(message, "Some message here")
            }
        },
		{
			Path: "delete",
			Action: func (message *nats.Msg) {
				nachos.Reply(message, struct { Message string `json:"message"` }{Message:"Acknowledged"})
            }
        }
    }
}
nachos.Subscribe(route)
```

You can add before and end actions by simply adding them to the route declaration, if the route has children nodes then 
all the before and end actions will be inherited by the children nodes. An example of a before action would be:
```go
route := nachos.Route{
	Path: "quest.delete",
	BeforeAction: []BeforeActionEvent{
	    func(message *nats.Msg) bool {
			if message.Header.Get("Magic-Word") !== "Please" {
				nachos.Reply(message, "Invalid magic word, please try better!")
				return false
            }
			return true
        },   	
    }
	Action: func (message *nats.Msg) {
        nachos.Reply(message, "Yay, you sent the right magic word!")
    }
}
```