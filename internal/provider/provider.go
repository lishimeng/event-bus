package provider

import (
	"gitee.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/go-log"
)

type MessageListener func(m message.Message)

type Provider interface {
	OnMessage(m *message.Message)   //
	Publish(m message.Message)      // 发布
	Subscribe(ch message.Channel)   // 订阅(一次性)
	UnSubscribe(ch message.Channel) // 取消订阅(一次性)
	SetMessageListener(MessageListener)
}

// RespListener 回调结果
type RespListener func(m message.Message)

type BaseProvider struct {
	decodeHandlers []MessageHandler
	encodeHandlers []MessageHandler
}

func (b *BaseProvider) AddDecodeHandler(handler MessageHandler) {
	b.decodeHandlers = append(b.decodeHandlers, handler)
}
func (b *BaseProvider) AddEncodeHandler(handler MessageHandler) {
	b.encodeHandlers = append(b.decodeHandlers, handler)
}

func (b *BaseProvider) PrePublish(m *message.Message) (err error) {
	log.Info("pre publish[encode]")
	ctx := make(map[string]any)
	for _, handler := range b.encodeHandlers {
		err = handler(m, ctx)
		if err != nil {
			return
		}
	}
	return
}

func (b *BaseProvider) OnMessage(m *message.Message) {
	log.Info("handleMessage: %s[%s]<-%s", m.RequestId, m.ReferId, m.Route)
	var err error
	ctx := make(map[string]any)
	for _, handler := range b.decodeHandlers {
		err = handler(m, ctx)
		if err != nil {
			break
		}
	}
	return
}
