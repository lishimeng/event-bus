package msgRecord

import (
	"context"

	"github.com/lishimeng/go-log"
)

type MsgRecord struct {
	Id      string
	Topic   string
	Payload string
}

var buffer chan MsgRecord
var enable = false

func OnMessage(rmqMsgId string, topic string, payload string) {
	if !enable {
		return
	}
	var r MsgRecord
	r.Id = rmqMsgId
	r.Topic = topic
	r.Payload = payload
	buffer <- r
}

func StartSave(ctx context.Context) {

	enable = true
	for {
		select {
		case <-ctx.Done():
			return
		case r := <-buffer:
			save(r)
		}
	}
}

func save(r MsgRecord) {
	defer func() {
		if r := recover(); r != nil {
			log.Info("Recovered in save %v", r)
		}
	}()
	//TODO
}
