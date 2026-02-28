package proxy

import (
	"context"

	"gitee.com/lishimeng/event-bus/providers/RocketMqProvider/msgRecord"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/lishimeng/go-log"
)

type Publisher struct {
	ctx          context.Context
	conf         Conf // 连接配置
	sendEnable   bool
	msgBuf       chan *rmq.Message
	topics       []string
	messageGroup string
}

func (p *Publisher) start() {
	log.Info("init publisher")
	p.msgBuf = make(chan *rmq.Message, 256)
	p.forever()
}

func (p *Publisher) forever() {
	log.Info("start publisher process")
	defer log.Info("end publisher process")
	for {
		select {
		case <-p.ctx.Done():
			return
		default:
			p.work()
		}
	}
}

func (p *Publisher) work() {
	log.Info("start publisher")
	defer func() {
		if err := recover(); err != nil {
			log.Info("publisher work panic")
			log.Info(err)
			return
		}
	}()
	producer, err := rmq.NewProducer(&rmq.Config{
		Endpoint: p.conf.Host,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    p.conf.AppId,
			AccessSecret: p.conf.Secret,
		},
	},
		rmq.WithTopics(p.topics...),
	)
	if err != nil {
		log.Info("connect endpoint fail")
		log.Info(err)
		return
	}
	// start producer
	err = producer.Start()
	if err != nil {
		log.Info("start publisher fail")
		log.Info(err)
		return
	}
	p.sendEnable = true
	// graceful stop producer
	defer func() {
		_ = producer.GracefulStop()
	}()
	log.Info("start publisher success, wait for exit")
	for {
		select {
		case <-p.ctx.Done():
			return
		case msg := <-p.msgBuf:
			var resp []*rmq.SendReceipt
			resp, err = producer.Send(p.ctx, msg)
			if err != nil {
				log.Info(err)
			}
			for i := 0; i < len(resp); i++ {
				var r = resp[i]
				log.Debug(">>> send message id: %s\n", r.MessageID)
				msgRecord.OnMessage(r.MessageID, msg.Topic, string(msg.Body))
			}
		}
	}
}

func (p *Publisher) Publish(topic string, payload []byte) {
	log.Debug("publish to topic:%s", topic)
	if !p.sendEnable {
		return
	}
	msg := &rmq.Message{
		Topic: topic,
		Body:  payload,
	}
	msg.SetMessageGroup(p.messageGroup)
	p.msgBuf <- msg
	log.Debug("publish msg to buffer")
	return
}
