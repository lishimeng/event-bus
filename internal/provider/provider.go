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
	msgListener    CallbackListener
	decodeHandlers []MessageHandler
	encodeHandlers []MessageHandler
}

func (b *BaseProvider) AddDecodeHandler(handler MessageHandler) {
	b.decodeHandlers = append(b.decodeHandlers, handler)
}
func (b *BaseProvider) AddEncodeHandler(handler MessageHandler) {
	b.encodeHandlers = append(b.decodeHandlers, handler)
}

func (b *BaseProvider) SetMsgListener(listener CallbackListener) {
	b.msgListener = listener
}
func (b *BaseProvider) SetRespListener(listener CallbackListener) {
	b.msgListener = listener
}

func (b *BaseProvider) PrePublish(m message.Message) (err error) {
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

func (b *BaseProvider) OnMessage(m message.Message) {
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
