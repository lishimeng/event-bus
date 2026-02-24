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

type Client struct {
	ctx               context.Context
	config            Conf
	lock              sync.RWMutex
	closeMap          map[string]func()
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

func (c *Client) initClient() {
	c.closeMap = make(map[string]func())
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
	go func() {
		c.publisher.start()
	}()
}

func (c *Client) Publish(topic string, data []byte) {
	c.publisher.Publish(topic, data)
}

func (c *Client) Subscribe(topic string, group string, fn OnMessageFunc) {
	go c.subscribe(topic, group, fn)
}

// Subscribe 阻塞
func (c *Client) subscribe(topic string, group string, fn OnMessageFunc) {
	ctx, cancel := context.WithCancel(c.ctx)
	c.closeMap[topic] = cancel
	defer cancel()
	var subscriber = Subscriber{
		ctx:           ctx,
		topic:         topic,
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
	c.close(topic)
}

func (c *Client) closeAll() {
	var list = c.genList()
	for _, name := range list {
		c.close(name)
	}
}

func (c *Client) genList() (list []string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for name := range c.closeMap {
		list = append(list, name)
	}
	return
}

func (c *Client) close(name string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	stopFunc, ok := c.closeMap[name]
	if !ok {
		return
	}
	stopFunc()
	delete(c.closeMap, name)
}
