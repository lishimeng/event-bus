package provider

import (
	"errors"
	"fmt"

	"github.com/lishimeng/event-bus/internal/channel"
	"github.com/lishimeng/event-bus/internal/db"
	"github.com/lishimeng/event-bus/internal/message"
)

var ChannelChkHandler = func(category db.RouteCategory) MessageHandler {
	return func(msg *message.Message, ctx map[string]any) (err error) {
		ch, err := channel.GetChannel(msg.Route, category)
		if err != nil {
			return
		}
		if ch.Category != category { // TODO 改进和这个判断不需要了, get_channel已经判断
			err = errors.New(fmt.Sprintf("[subscriber:%s]channel doesn't exist", msg.Route))
			return
		}
		return
	}
}
