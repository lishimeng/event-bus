package provider

import (
	"github.com/lishimeng/event-bus/internal/channel"
	"github.com/lishimeng/event-bus/internal/db"
	"github.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/go-log"
)

// TlsDecryptHandler 解密
var TlsDecryptHandler MessageHandler = func(msg *message.Message, ctx map[string]any) (err error) {
	ch, err := channel.GetChannel(msg.Route, db.Subscribe) // 订阅类型通道同来解密数据(接收时解密)
	if err != nil {
		return
	}
	biz, err := message.Decrypt(msg.Payload, ch)
	if err != nil {
		log.Info("解密流程失败")
		log.Info(err)
		return
	}
	ctx["biz"] = biz
	msg.Biz = biz
	return
}
