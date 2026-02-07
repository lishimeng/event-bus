package proc

import (
	"context"

	"gitee.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/go-log"
)

type Engine struct {
}

var EngineInstance *Engine

// Subscribe 订阅
func (h *Engine) Subscribe(ch message.Channel) {
	log.Info("subscribe %s[%s] to %s", ch.Code, ch.Name, ch.Route)
	// TODO
}

// Unsubscribe 反订阅
func (h *Engine) Unsubscribe(ch message.Channel) {

}

func (h *Engine) OnMessage(m message.Message) {

}

// Publish 发布消息
func (h *Engine) Publish(msg message.Message) {

}

func Start(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		one()
	}
}

func one() {
	//
	var m message.Message
	Callback(m)
}
