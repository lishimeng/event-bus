package RocketMqProvider

import "gitee.com/lishimeng/event-bus/internal/provider"

type RocketMqProvider struct {
	provider.BaseProvider
}

func (p *RocketMqProvider) Init(b provider.BaseProvider) {
	p.BaseProvider = b

	p.AddHandler(provider.DataRecordMsgHandler)
	p.AddHandler(provider.ChannelChkHandler)
	p.AddHandler(provider.TlsDecryptHandler)
	p.AddHandler(msgExecHandler)
}
