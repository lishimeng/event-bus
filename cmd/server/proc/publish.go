package proc

import (
	"encoding/json"

	"gitee.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/go-log"
)

func Publish(m message.Message) {
	// TODO
	log.Info("publish message")
	bs, _ := json.Marshal(m)
	log.Info(string(bs))
	if EngineInstance == nil {
		return
	}
	EngineInstance.Publish(m)
}
