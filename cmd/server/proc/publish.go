package proc

import (
	"gitee.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/go-log"
)

func Publish(m message.Message) {
	// TODO
	log.Info("publish message")
	if EngineInstance == nil {
		return
	}
	EngineInstance.Publish(m)
}
