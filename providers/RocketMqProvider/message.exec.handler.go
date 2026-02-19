package RocketMqProvider

import (
	"errors"

	"gitee.com/lishimeng/event-bus/internal/channel"
	"gitee.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/go-log"
)

var msgExecHandler = func(msg *message.Message, ctx map[string]any) (err error) {
	biz, ok := ctx["biz"]
	if !ok {
		err = errors.New("biz not exist")
		return
	}
	var bizMessage message.BizMessage
	bizMessage, ok = biz.(message.BizMessage)
	if !ok {
		err = errors.New("bizMessage not exist")
		return
	}
	log.Info("biz_msg: %s[%s]", bizMessage.Action, bizMessage.Method)
	ch, err := channel.GetChannel(msg.Route)
	if err != nil {
		return
	}
	log.Info("callback_uri: %s", ch.Callback)
	// TODO 执行
	return
}
