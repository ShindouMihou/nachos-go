package nachos

import (
	"errors"
	"github.com/nats-io/nats.go"
)

func Attach(conn *nats.Conn) error {
	for path, route := range Routes {
		subscription, err := conn.Subscribe(path, func(message *nats.Msg) {
			Handle(route, message)
		})
		if err != nil {
			return errors.Join(errors.New("failed to subscribe to "+path+": "), err)
		}
		Subscriptions[path] = subscription
	}
	return nil
}

func Handle(route *Route, message *nats.Msg) {
	go func() {
		canContinue := true
		for _, beforeAction := range route.BeforeAction {
			next := beforeAction(message)
			if !next {
				canContinue = false
				break
			}
		}
		if !canContinue {
			return
		}
		route.Action(message)
		for _, endAction := range route.EndAction {
			go endAction(message)
		}
	}()
}
