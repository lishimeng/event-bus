package provider

import (
	"gitee.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/go-log"
)

type Provider interface {
	OnMessage(m message.Message)
	Publish(m message.Message)
	Subscribe(ch message.Channel)
}

// RespListener 回调结果
type RespListener func(m message.Message)

// CallbackListener 回调
type CallbackListener func(biz message.BizMessage) (resp map[string]any, err error)
type BaseProvider struct {
	msgListener CallbackListener
	handlers    []MessageHandler
}

func (b *BaseProvider) AddHandler(handler MessageHandler) {
	b.handlers = append(b.handlers, handler)
}

func (b *BaseProvider) SetMsgListener(listener CallbackListener) {
	b.msgListener = listener
}
func (b *BaseProvider) SetRespListener(listener CallbackListener) {
	b.msgListener = listener
}

func (b *BaseProvider) OnMessage(m message.Message) {
	log.Info("OnMessage: %s[%s]<-%s", m.RequestId, m.ReferId, m.Route)
	var err error
	ctx := make(map[string]any)
	for _, handler := range b.handlers {
		err = handler(m, ctx)
		if err != nil {
			break
		}
	}
	return
}
