package proxy

import (
	"context"
	"sync"

	"github.com/lishimeng/go-log"
)

type OptFunc func(opt *Conf)

type Auth struct {
	AppId  string
	Secret string
}

type Conf struct {
	Host                  string
	publisherMessageGroup string
	publisherTopics       []string
	consumerGroup         string
	consumerTopics        []string
	consumerHandler       OnMessageFunc
	Auth
}

var WithEndpoint = func(endpoint string) OptFunc {
	return func(opt *Conf) {
		opt.Host = endpoint
	}
}
var WithAuth = func(appId, secret string) OptFunc {
	return func(opt *Conf) {
		opt.Auth.AppId = appId
		opt.Auth.Secret = secret
	}
}
var WithPublisherConfigs = func(messageGroup string, topics ...string) OptFunc {
	return func(opt *Conf) {
		opt.publisherMessageGroup = messageGroup
		opt.publisherTopics = topics
	}
}

var WithConsumerConfigs = func(messageGroup string, topics ...string) OptFunc {
	return func(opt *Conf) {
		opt.consumerGroup = messageGroup
		opt.consumerTopics = topics
	}
}

var WithConsumerHandler = func(handler OnMessageFunc) OptFunc {
	return func(opt *Conf) {
		opt.consumerHandler = handler
	}
}

type Client struct {
	ctx               context.Context
	config            Conf
	lock              sync.RWMutex
	globalMessageFunc OnMessageFunc

	publisher *Publisher
}

func New(ctx context.Context, opts ...OptFunc) (c *Client) {
	c = &Client{ctx: ctx}
	c.config = Conf{}
	for _, opt := range opts {
		opt(&c.config)
	}
	c.initClient()
	return
}

func (c *Client) Start() {
	go func() {
		c.publisher.start()
	}()
	go func() {
		c.subscribe(c.config.consumerTopics, c.config.consumerGroup, c.config.consumerHandler)
	}()
}

func (c *Client) initClient() {
	c.initPublisher()
}

func (c *Client) initPublisher() {
	log.Info("init publisher")
	c.publisher = &Publisher{
		ctx:          c.ctx,
		conf:         c.config,
		topics:       c.config.publisherTopics,
		messageGroup: c.config.publisherMessageGroup,
	}
}

func (c *Client) Publish(topic string, data []byte) {
	c.publisher.Publish(topic, data)
}

// Subscribe 阻塞
func (c *Client) subscribe(topics []string, group string, fn OnMessageFunc) {
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()
	var subscriber = Subscriber{
		ctx:           ctx,
		topics:        topics,
		conf:          c.config,
		consumerGroup: group,
		onMessage:     fn,
	}
	if subscriber.onMessage == nil { // 优先使用指定回调
		subscriber.onMessage = c.globalMessageFunc
	}
	subscriber.Subscribe()
}

func (c *Client) UnSubscribe(topic string) {

}
