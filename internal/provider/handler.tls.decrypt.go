package provider

import (
	"gitee.com/lishimeng/event-bus/internal/channel"
	"gitee.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/go-log"
)

// TlsDecryptHandler 解密
var TlsDecryptHandler MessageHandler = func(msg *message.Message, ctx map[string]any) (err error) {
	ch, err := channel.GetChannel(msg.Route)
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
