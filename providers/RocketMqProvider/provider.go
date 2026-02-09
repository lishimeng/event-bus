package RocketMqProvider

import (
	"gitee.com/lishimeng/event-bus/internal/message"
	"gitee.com/lishimeng/event-bus/internal/provider"
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

	p.AddHandler(provider.DataRecordMsgHandler) // 记录接收
	p.AddHandler(provider.ChannelChkHandler)    // 检查通道支持
	p.AddHandler(provider.TlsDecryptHandler)    // 解密数据
	p.AddHandler(msgExecHandler)                // 回调
}

func (p *RocketMqProvider) Publish(m message.Message) {

}

func (p *RocketMqProvider) Subscribe(ch message.Channel) {

}
