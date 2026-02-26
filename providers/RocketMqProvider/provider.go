package RocketMqProvider

import (
	"encoding/json"

	"gitee.com/lishimeng/event-bus/internal/db"
	"gitee.com/lishimeng/event-bus/internal/message"
	"gitee.com/lishimeng/event-bus/internal/provider"
	"gitee.com/lishimeng/event-bus/providers/RocketMqProvider/msgRecord"
	"gitee.com/lishimeng/event-bus/providers/RocketMqProvider/proxy"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/lishimeng/go-log"
)

type RocketMqProvider struct {
	provider.BaseProvider
	client          *proxy.Client
	rmqConfig       RmqConfig
	messageListener provider.MessageListener
}

func New(client *proxy.Client, cfg RmqConfig) (p provider.Provider) {
	h := &RocketMqProvider{
		client:    client,
		rmqConfig: cfg,
	}
	h.Init()
	p = h
	return
}

func (p *RocketMqProvider) Init() {
	// ----------subscribe-----------------------------------------
	p.AddDecodeHandler(provider.DataRecordMsgHandler)             // 记录接收
	p.AddDecodeHandler(provider.ChannelChkHandler(db.Subscriber)) // 检查通道支持
	p.AddDecodeHandler(provider.TlsDecryptHandler)                // 解密数据
	//p.AddDecodeHandler(msgExecHandler)                            // 回调

	// ---------publish--------------------------------------------------
	p.AddEncodeHandler(provider.DataRecordMsgHandler)
	p.AddEncodeHandler(provider.ChannelChkHandler(db.Publish))
	p.AddEncodeHandler(provider.TlsEncryptHandler) // 加密
	p.AddEncodeHandler(rmqMsgPubHandler)
}

func (p *RocketMqProvider) Publish(m message.Message) {

	err := p.BaseProvider.PrePublish(&m)
	if err != nil {
		log.Info("publish fail")
		log.Info(err)
		return
	}
	log.Info("publish success")

}

func (p *RocketMqProvider) Subscribe(ch message.Channel) {

	topic := ch.Route
	subCfg, err := p.rmqConfig.GetSubscriber(topic)
	if err != nil {
		log.Info("subscribe fail, topic not supported[%s]", topic)
		log.Info(err)
		return
	}
	p.client.Subscribe(subCfg.Topic, subCfg.ConsumerGroup, func(mv *rmq.MessageView) {
		msgRecord.OnMessage(mv.GetMessageId(), mv.GetTopic(), string(mv.GetBody()))
		var m message.Message // TODO

		err = json.Unmarshal(mv.GetBody(), &m)
		if err != nil {
			log.Info(err)
			return
		}
		m.Route = mv.GetTopic() // 修正一次router,与rmq中一致
		p.onMessage(m)
	})
}

func (p *RocketMqProvider) UnSubscribe(ch message.Channel) {
	key := ch.Route
	p.client.UnSubscribe(key)
}

func (p *RocketMqProvider) onMessage(m message.Message) {
	p.BaseProvider.OnMessage(&m)
	if p.messageListener != nil {
		p.messageListener(m)
	}
}

func (p *RocketMqProvider) SetMessageListener(listener provider.MessageListener) {
	p.messageListener = listener
}
