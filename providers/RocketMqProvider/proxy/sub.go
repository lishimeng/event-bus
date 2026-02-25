package proxy

import (
	"context"
	"time"

	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/lishimeng/go-log"
)

type OnMessageFunc func(message *rmq.MessageView)

const (
	awaitDuration = time.Second * 5
)

var (
	// maximum waiting time for receive func
	// maximum number of messages received at one time
	maxMessageNum int32 = 16
	// invisibleDuration should > 20s
	invisibleDuration = time.Second * 20
	// receive concurrency
	receiveConcurrency = 1
)

type Subscriber struct {
	ctx           context.Context
	conf          Conf          // 连接配置
	topic         string        // 订阅
	consumerGroup string        // 消费组
	onMessage     OnMessageFunc // 回调
}

// 不间断运行, 监听ctx结束
func (sub *Subscriber) runForever() {
	//var retry = 0
	for {
		select {
		case <-sub.ctx.Done():
			return
		default:
			_, _ = sub.createConsumer()

		}
	}
}

// 一次完整运行, 遇到不可预知问题就结束, 但保证线程safety
func (sub *Subscriber) runOnce(consumer rmq.SimpleConsumer) (err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Info(r)
		}
		log.Info("stop receive message[%s<-(%s)]", sub.topic, sub.consumerGroup)
		_ = consumer.GracefulStop()
	}()

	log.Info("start receive message")
	for {
		select {
		case <-sub.ctx.Done():
			return
		default:
			err = sub.workOnce(consumer)
			if err != nil {
				log.Info(err)
				time.Sleep(time.Second * 1)
			}
		}
	}
}

func (sub *Subscriber) createConsumer() (consumer rmq.PushConsumer, err error) {
	consumer, err = rmq.NewPushConsumer(&rmq.Config{
		Endpoint:      sub.conf.Host,
		ConsumerGroup: sub.consumerGroup,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    sub.conf.AppId,
			AccessSecret: sub.conf.Secret,
		},
	},

		rmq.WithPushAwaitDuration(awaitDuration),
		rmq.WithPushSubscriptionExpressions(map[string]*rmq.FilterExpression{
			sub.topic: rmq.SUB_ALL,
		}),
		rmq.WithPushMessageListener(&rmq.FuncMessageListener{
			Consume: func(mv *rmq.MessageView) rmq.ConsumerResult {
				sub.handleMessage(mv)
				return rmq.SUCCESS
			},
		}),
		rmq.WithPushConsumptionThreadCount(20),
		rmq.WithPushMaxCacheMessageCount(1024),
	)
	if err != nil {
		log.Info(err)
		return
	}
	err = consumer.Start()
	if err != nil {
		log.Info(err)
		return
	}
	//err = consumer.Subscribe(sub.topic, rmq.SUB_ALL)
	return
}

func (sub *Subscriber) Subscribe() {
	log.Info("subscribe topic:%s[%s]", sub.topic, sub.consumerGroup)
	sub.runForever()
	return
}

func (sub *Subscriber) workOnce(consumer rmq.SimpleConsumer) (err error) {
	mvs, err := consumer.Receive(context.TODO(), maxMessageNum, invisibleDuration) // 轮询方式
	if err != nil {
		rpcStatus, ok := rmq.AsErrRpcStatus(err)
		if ok {
			if rpcStatus.Code == 40401 {
				log.Debug("%s", rpcStatus.Error())
				err = nil
				return
			}
		}
		log.Info("receive message error: %s", err.Error())
		return
	}
	// ack message
	for _, mv := range mvs {
		sub.handleMessage(mv)
		if err = consumer.Ack(context.TODO(), mv); err != nil {
			log.Info("message ack error: %s[%s]"+mv.GetMessageId(), err.Error())
		} else {
			log.Info("message ack success: %s", mv.GetMessageId())
		}
	}
	return
}

func (sub *Subscriber) handleMessage(mv *rmq.MessageView) {
	log.Debug("<<< receive message: [%s]", mv.GetMessageId())
	if sub.onMessage == nil {
		return
	}
	sub.onMessage(mv)
}
