package RocketMqProvider

import (
	"gitee.com/lishimeng/event-bus/internal/db"
	"gitee.com/lishimeng/event-bus/internal/message"
	"gitee.com/lishimeng/event-bus/internal/provider"
	"github.com/lishimeng/go-log"
)

type RocketMqProvider struct {
	provider.BaseProvider
}

func New(base provider.BaseProvider) (p provider.Provider) {
	h := &RocketMqProvider{
		BaseProvider: base,
	}
	p = h
	return
}

func (p *RocketMqProvider) Init(b provider.BaseProvider) {
	p.BaseProvider = b

	// ----------subscribe-----------------------------------------
	p.AddDecodeHandler(provider.DataRecordMsgHandler)             // 记录接收
	p.AddDecodeHandler(provider.ChannelChkHandler(db.Subscriber)) // 检查通道支持
	p.AddDecodeHandler(provider.TlsDecryptHandler)                // 解密数据
	p.AddDecodeHandler(msgExecHandler)                            // 回调

	// ---------publish--------------------------------------------------
	p.AddEncodeHandler(provider.DataRecordMsgHandler)
	p.AddEncodeHandler(provider.ChannelChkHandler(db.Publish))
	p.AddEncodeHandler(provider.TlsEncryptHandler) // 加密
}

func (p *RocketMqProvider) Publish(m message.Message) {

	err := p.BaseProvider.PrePublish(m)
	if err != nil {
		log.Info("pre publish fail")
		return
	}

	// TODO publish
}

func (p *RocketMqProvider) Subscribe(ch message.Channel) {

	var m message.Message // TODO
	p.onMessage(m)
}

func (p *RocketMqProvider) onMessage(m message.Message) {
	p.BaseProvider.OnMessage(m)
}
