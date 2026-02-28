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
	var retry = 0
	const maxRetryTimes = 90
	for {
		select {
		case <-sub.ctx.Done():
			return
		default:
			consumer, err := sub.createConsumer()
			if err != nil {
				log.Info(err)
				retry++
				if retry > maxRetryTimes {
					retry = maxRetryTimes
				}
				time.Sleep(time.Second * time.Duration(retry))
			} else {
				retry = 0
				log.Info("consumer loop start[retry:%d]", retry)
				err = sub.consumerLoop(consumer)
				if err != nil {
					log.Info("consumer start fail")
					log.Info(err)
				} else {
					log.Info("consumer loop end")
				}
			}
		}
	}
}

func (sub *Subscriber) consumerLoop(consumer rmq.PushConsumer) (err error) {
	log.Info("start consumer loop %s", sub.topic)
	err = consumer.Start()
	if err != nil {
		log.Info("consumer start fail")
		log.Info(err)
		return
	}
	defer func() {
		_ = consumer.GracefulStop()
	}()

	log.Info("consumer loop start success, wait for exit")
	select {
	case <-sub.ctx.Done():
		log.Info("consumer loop exit")
		return
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
	return
}

func (sub *Subscriber) Subscribe() {
	log.Info("subscribe topic:%s[%s]", sub.topic, sub.consumerGroup)
	sub.runForever()
	return
}

func (sub *Subscriber) handleMessage(mv *rmq.MessageView) {
	log.Debug("<<< receive message: [%s]", mv.GetMessageId())
	if sub.onMessage == nil {
		return
	}
	sub.onMessage(mv)
}
