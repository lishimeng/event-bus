package proc

import (
	"context"

	"github.com/lishimeng/event-bus/internal/channel"
	"github.com/lishimeng/event-bus/internal/db"
	"github.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/event-bus/internal/provider"
	"github.com/lishimeng/go-log"
)

type Engine struct {
	worker provider.Provider
}

var EngineInstance *Engine

func NewEngine(worker provider.Provider) *Engine {
	return &Engine{worker: worker}
}

// Subscribe 订阅
func (h *Engine) Subscribe(ch message.Channel) {
	log.Info("subscribe %s[%s] to %s", ch.Code, ch.Name, ch.Route)
	h.worker.Subscribe(ch)
}

// Unsubscribe 反订阅
func (h *Engine) Unsubscribe(ch message.Channel) {
	log.Info("unsubscribe %s[%s] to %s", ch.Code, ch.Name, ch.Route)
	h.worker.UnSubscribe(ch)
}

func (h *Engine) OnMessage(m message.Message) {
	bizMessage := m.Biz

	log.Info("biz_msg: %s[%s]", bizMessage.Action, bizMessage.Method)
	ch, err := channel.GetChannel(m.Route, db.Subscribe) // 订阅类别(业务系统订阅)
	if err != nil {
		return
	}
	log.Info("callback_uri: %s", ch.Callback)
	err = Callback(m)
	if err != nil {
		log.Info(err)
	}
}

// Publish 发布消息
func (h *Engine) Publish(msg message.Message) {
	h.worker.Publish(msg)
}

func Start(ctx context.Context) {

}
