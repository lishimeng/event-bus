package RocketMqProvider

import "gitee.com/lishimeng/event-bus/internal/provider"

type RocketMqProvider struct {
	provider.BaseProvider
}

func (p *RocketMqProvider) Init(b provider.BaseProvider) {
	p.BaseProvider = b

	p.AddHandler(provider.DataRecordMsgHandler) // 记录接收
	p.AddHandler(provider.ChannelChkHandler)    // 检查通道支持
	p.AddHandler(provider.TlsDecryptHandler)    // 解密数据
	p.AddHandler(msgExecHandler)                // 回调
}
