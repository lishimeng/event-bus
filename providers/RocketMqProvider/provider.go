package RocketMqProvider

import (
	"context"
	"encoding/json"

	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/lishimeng/event-bus/internal/db"
	"github.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/event-bus/internal/provider"
	"github.com/lishimeng/event-bus/providers/RocketMqProvider/msgRecord"
	"github.com/lishimeng/event-bus/providers/RocketMqProvider/proxy"
	"github.com/lishimeng/go-log"
	"github.com/lishimeng/x/container"
)

type RocketMqProvider struct {
	ctx context.Context
	provider.BaseProvider
	client          *proxy.Client
	rmqConfig       RmqConfig
	publishTopics   []string
	subscribeTopics []string
	messageListener provider.MessageListener
}

func New(ctx context.Context, cfg RmqConfig, publishTopics []string, subscribeTopics []string) (p provider.Provider) {
	h := &RocketMqProvider{
		ctx:             ctx,
		rmqConfig:       cfg,
		publishTopics:   publishTopics,
		subscribeTopics: subscribeTopics,
	}
	h.Init()
	p = h
	go func() {
		h.client.Start()
	}()
	return
}

func (p *RocketMqProvider) createProxy() {
	var cfg = p.rmqConfig
	var publishTopics = p.publishTopics
	var subscribeTopics = p.subscribeTopics
	client := proxy.New(p.ctx,
		proxy.WithEndpoint(p.rmqConfig.Endpoint),
		proxy.WithAuth(cfg.AppId, cfg.Secret),
		proxy.WithPublisherConfigs(cfg.Publisher.MessageGroup, publishTopics...),
		proxy.WithConsumerConfigs(cfg.Subscribers[0].ConsumerGroup, subscribeTopics...),
		proxy.WithConsumerHandler(func(mv *rmq.MessageView) {
			msgRecord.OnMessage(mv.GetMessageId(), mv.GetTopic(), string(mv.GetBody()))
			var m message.Message // TODO

			err := json.Unmarshal(mv.GetBody(), &m)
			if err != nil {
				log.Info(err)
				return
			}
			m.Route = mv.GetTopic() // 修正一次router,与rmq中一致
			p.onMessage(m)
		}),
	)
	container.Add(&client)
	p.client = client
}

func (p *RocketMqProvider) Init() {
	// ----------subscribe-----------------------------------------
	p.AddDecodeHandler(provider.DataRecordMsgHandler)            // 记录接收
	p.AddDecodeHandler(provider.ChannelChkHandler(db.Subscribe)) // 检查通道支持
	p.AddDecodeHandler(provider.TlsDecryptHandler)               // 解密数据
	//p.AddDecodeHandler(msgExecHandler)                            // 回调

	// ---------publish--------------------------------------------------
	p.AddEncodeHandler(provider.DataRecordMsgHandler)
	p.AddEncodeHandler(provider.ChannelChkHandler(db.PublishTo))
	p.AddEncodeHandler(provider.TlsEncryptHandler) // 加密
	p.AddEncodeHandler(rmqMsgPubHandler)

	p.createProxy()
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
	log.Info("dummy subscribe: %s", ch.Name)
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
