package nachos

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
)

func marshal(message *nats.Msg, msg any) ([]byte, bool) {
	if contents, ok := msg.(string); ok {
		return []byte(contents), true
	}
	contents, err := json.Marshal(msg)
	if err != nil {
		Logger(Error, "Failed to marshal response to ", message.Subject, ": ", err)
		return nil, false
	}
	return contents, true
}

func Reply(message *nats.Msg, msg any) {
	if contents, valid := marshal(message, msg); valid {
		err := message.Respond(contents)
		if err != nil {
			Logger(Error, "Failed to send response to ", message.Subject, ": ", err)
			return
		}
	}
}

func ReplyMsg(message *nats.Msg, base *nats.Msg, msg *any) {
	if contents, valid := marshal(message, msg); valid {
		base.Data = contents
		err := message.RespondMsg(base)
		if err != nil {
			Logger(Error, "Failed to send response to ", message.Subject, ": ", err)
			return
		}
	}
}
