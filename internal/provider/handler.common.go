package provider

import (
	"errors"
	"fmt"
	"time"

	"gitee.com/lishimeng/event-bus/internal/channel"
	"gitee.com/lishimeng/event-bus/internal/db"
	"gitee.com/lishimeng/event-bus/internal/message"
)

type MessageHandler func(msg message.Message, ctx map[string]any) (err error)

const (
	MaxMessageBuffer = 32
)

var ChannelChkHandler MessageHandler = func(msg message.Message, ctx map[string]any) (err error) {
	ch, err := channel.GetChannel(msg.Route)
	if err != nil {
		return
	}
	if ch.Category != db.Subscriber {
		err = errors.New(fmt.Sprintf("[subscriber:%s]channel doesn't exist", msg.Route))
		return
	}
	return
}

var queue = make(chan message.Message, MaxMessageBuffer)

// DataRecordMsgHandler 存储
var DataRecordMsgHandler MessageHandler = func(msg message.Message, ctx map[string]any) (err error) {
	select {
	case queue <- msg:
	case <-time.After(time.Millisecond * 10):
		return
	}
	return
}
