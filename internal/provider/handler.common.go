package provider

import (
	"time"

	"gitee.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/go-log"
)

type MessageHandler func(msg message.Message, ctx map[string]any) (err error)

const (
	MaxMessageBuffer = 32
)

var queue = make(chan message.Message, MaxMessageBuffer)

// DataRecordMsgHandler 存储
var DataRecordMsgHandler MessageHandler = func(msg message.Message, ctx map[string]any) (err error) {
	select {
	case queue <- msg:
		log.Info("save msg to queue:%s", msg.Route)
	case <-time.After(time.Millisecond * 10):
		return
	}
	return
}
