package nachos

import (
	"errors"
	"github.com/nats-io/nats.go"
)

var DefaultQueueGroup *QueueGroup = nil

func Attach(conn *nats.Conn) error {
	for path, route := range Routes {
		var subscription *nats.Subscription
		SubscribeRegularly := func(path string, route *Route) error {
			s, err := conn.Subscribe(path, func(message *nats.Msg) {
				Handle(route, message)
			})
			if err != nil {
				return errors.Join(errors.New("failed to subscribe to "+path+": "), err)
			}
			subscription = s
			return nil
		}
		if route.QueueGroup != nil || DefaultQueueGroup != nil {
			var queueGroup = route.QueueGroup
			if queueGroup == nil {
				queueGroup = DefaultQueueGroup
			}
			if queueGroup.Enabled {
				s, err := conn.QueueSubscribe(path, queueGroup.Name, func(message *nats.Msg) {
					Handle(&route, message)
				})
				if err != nil {
					return errors.Join(errors.New("failed to subscribe to "+path+": "), err)
				}
				subscription = s
			} else {
				if err := SubscribeRegularly(path, &route); err != nil {
					return err
				}
			}
		} else {
			if err := SubscribeRegularly(path, &route); err != nil {
				return err
			}
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
