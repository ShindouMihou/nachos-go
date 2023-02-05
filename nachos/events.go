package nachos

import "github.com/nats-io/nats.go"

type BeforeActionEvent = func(message *nats.Msg) bool
type ActionEvent = func(message *nats.Msg)
type EndActionEvent = func(message *nats.Msg)
