package RocketMqProvider

import (
	"encoding/json"

	"gitee.com/lishimeng/event-bus/internal/message"
	"gitee.com/lishimeng/event-bus/providers/RocketMqProvider/proxy"
	"github.com/lishimeng/go-log"
	"github.com/lishimeng/x/container"
)

var rmqMsgPubHandler = func(msg *message.Message, ctx map[string]any) (err error) {

	var rmqProxy *proxy.Client
	err = container.Get(&rmqProxy)
	if err != nil {
		log.Info("rmq proxy not exist")
		return err
	}
	route := msg.Route
	bs, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	rmqProxy.Publish(route, bs)
	return
}
