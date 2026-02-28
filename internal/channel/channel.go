package channel

import (
	"crypto/rsa"
	"errors"

	"gitee.com/lishimeng/event-bus/internal/db"
	"gitee.com/lishimeng/event-bus/internal/message"
	"gitee.com/lishimeng/event-bus/internal/tls/cypher"
	"github.com/lishimeng/go-log"
)

type Manager struct {
	// channel列表: Category_router-> channel实例
	channels map[string]message.Channel
}

func (m *Manager) Register(c message.Channel) {
	m.channels[c.GetKey()] = c
}

func (m *Manager) Get(key string) (ch message.Channel, err error) {
	ch, ok := m.channels[key]
	if !ok {
		err = errors.New("channel not found")
	}
	return
}

func (m *Manager) GetCh(route string, c db.RouteCategory) (ch message.Channel, err error) {
	key := message.GenKey(c, route)
	return m.Get(key)
}

var managerSingleton *Manager

func init() {
	managerSingleton = &Manager{channels: make(map[string]message.Channel)}
}

func LoadChannel(config db.ChannelConfig) (ch message.Channel, err error) {
	ch.Code = config.Code
	ch.Name = config.Name
	ch.Category = config.Category
	ch.Route = config.Router
	ch.Callback = config.Callback
	if config.UseSecurity == 1 {
		ch.UseTls = true
		ch.Cipher, err = resolveChSecret(config.Security, config.Category)
		if err != nil {
			return
		}
		if config.Category == db.Subscribe { // subscriber预先创建session
			err = ch.RefreshSession() // 保证session不空
			if err != nil {
				return
			}
		}

	}

	// 全局通道
	managerSingleton.Register(ch)
	log.Info("load channel success. %s[%s]->%s:category:%s", ch.Code, ch.Name, ch.Route, ch.Category.String())

	// 分组通道
	switch ch.Category {
	case db.PublishTo:
		log.Info("publish channel register")
		publishers[ch.Route] = ch
	case db.Subscribe:
		log.Info("subscriber channel register")
		subscribers[ch.Route] = ch
	default:
		log.Info("not support channel type, %d", ch.Category)
		err = errors.New("not support channel type")
	}
	return
}

func resolveChSecret(s string, category db.RouteCategory) (c message.ChannelCipher, err error) {
	var chSecret db.ChannelSecurity
	err = chSecret.Unmarshal(s)
	if err != nil {
		log.Info("unmarshal channel secret")
		log.Info(err)
		return
	}

	var pubKey *rsa.PublicKey
	var priKey *rsa.PrivateKey
	if len(chSecret.RsaPem) > 0 {
		pubKey, err = cypher.LoadPublicKey([]byte(chSecret.RsaPem))
		if err != nil {
			return
		}
		c.RsaPubKey = pubKey // 通道中只加载公钥
	}
	if len(chSecret.RsaKey) > 0 {
		priKey, err = cypher.LoadPrivateKey([]byte(chSecret.RsaKey))
		if err != nil {
			return
		}
		c.RsaPriKey = priKey
	}

	return
}

// GetChannel 查询支持的通道(全局通道)
func GetChannel(route string, category db.RouteCategory) (ch message.Channel, err error) {
	ch, err = managerSingleton.GetCh(route, category)
	return
}
