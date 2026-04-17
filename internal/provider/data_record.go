package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lishimeng/app-starter"
	"github.com/lishimeng/event-bus/internal/db"
	"github.com/lishimeng/event-bus/internal/id"
	"github.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/event-bus/sdk"
	"github.com/lishimeng/go-log"
)

const recordQueueCap = 256

type recordJob struct {
	dir db.RouteCategory
	msg message.Message
}

var recordQueue = make(chan recordJob, recordQueueCap)

// StartRecordPersistence 在 DB 就绪后调用：异步将消息写入 data_record。
func StartRecordPersistence(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case job := <-recordQueue:
				insertDataRecord(job.dir, &job.msg)
			}
		}
	}()
}

func enqueueDataRecord(dir db.RouteCategory, msg *message.Message) {
	select {
	case recordQueue <- recordJob{dir: dir, msg: *msg}:
	default:
		log.Info("data_record queue full, drop route=%s", msg.Route)
	}
}

// DataRecordPublishHandler 出站（加密发送前）落库。
var DataRecordPublishHandler MessageHandler = func(msg *message.Message, ctx map[string]any) (err error) {
	if msg.Biz.Method == "" {
		if v, ok := ctx["biz"]; ok {
			if biz, ok := v.(sdk.BizMessage); ok {
				msg.Biz = biz
			}
		}
	}
	enqueueDataRecord(db.PublishTo, msg)
	return
}

// DataRecordSubscribeHandler 入站（从 MQ 消费解码链）落库。
var DataRecordSubscribeHandler MessageHandler = func(msg *message.Message, _ map[string]any) (err error) {
	enqueueDataRecord(db.Subscribe, msg)
	return
}

func insertDataRecord(dir db.RouteCategory, msg *message.Message) {
	defer func() {
		if r := recover(); r != nil {
			log.Info("insertDataRecord panic: %v", r)
		}
	}()
	base := msg.RequestId
	if base == "" {
		base = id.GenId()
	}
	// 同一条 MQ 消息先 Publish 落库再 Subscribe 落库时 RequestId 相同，而 code 列 unique；
	// 用 route 方向（1=Subscribe / 2=PublishTo）后缀区分两行，便于对账同一 messageId。
	code := fmt.Sprintf("%s#%d", base, int(dir))
	pl, err := json.Marshal(msg.Payload)
	if err != nil {
		log.Info("data_record marshal payload: %v", err)
		return
	}
	biz, err := json.Marshal(msg.Biz)
	if err != nil {
		log.Info("data_record marshal biz: %v", err)
		return
	}
	source := msg.Source
	if source == "" {
		source = msg.Route
	}
	rec := db.DataRecord{
		Code:       code,
		ReferCode:  msg.ReferId,
		Source:     source,
		Route:      dir,
		Payload:    string(pl),
		BizPayload: string(biz),
	}
	rec.Status = 1
	_, err = app.GetOrm().Context.Insert(&rec)
	if err != nil {
		log.Info("data_record insert fail: %v", err)
		return
	}
	log.Info("data_record saved code=%s dir=%s", rec.Code, dir.String())
}
