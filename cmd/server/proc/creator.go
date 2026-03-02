package proc

import (
	"github.com/lishimeng/event-bus/internal/channel"
	"github.com/lishimeng/event-bus/internal/db"
	"github.com/lishimeng/event-bus/internal/id"
	"github.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/event-bus/sdk"
	"github.com/lishimeng/go-log"
)

type MessageCreateOpt struct {
	Id       string
	ParentId string
}

type MessageCreateFunc func(opt *MessageCreateOpt)

var WithId = func(id string) MessageCreateFunc {
	return func(opt *MessageCreateOpt) {
		opt.Id = id
	}
}

var WithParentId = func(id string) MessageCreateFunc {
	return func(opt *MessageCreateOpt) {
		opt.ParentId = id
	}
}

func Create(destination string, biz sdk.BizMessage, opts ...MessageCreateFunc) (m message.Message, err error) {
	// 消息创建业务
	var opt MessageCreateOpt
	for _, o := range opts {
		o(&opt)
	}
	m.RequestId = opt.Id
	m.ReferId = opt.ParentId
	m.Route = destination
	if len(m.RequestId) == 0 {
		m.RequestId = id.GenId() // 默认id策略
	}

	// get destination
	ch, err := channel.GetManager().GetCh(destination, db.PublishTo) // 发送到目的地(从发送channel列表中查找路由)
	if err != nil {
		log.Info("not found channel:%s[%s]", destination, db.PublishTo.String())
		return
	}

	m.Payload, err = message.Encrypt(biz, ch)
	if err != nil {
		log.Info("encrypt fail")
		log.Info(err)
		return
	}

	return
}
