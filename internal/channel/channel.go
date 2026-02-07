package channel

import (
	"crypto/rsa"
	"encoding/json"
	"errors"

	"gitee.com/lishimeng/event-bus/internal/db"
	"gitee.com/lishimeng/event-bus/internal/message"
	"gitee.com/lishimeng/event-bus/internal/tls/cypher"
)

var channels map[string]message.Channel // route:channel

func init() {
	channels = make(map[string]message.Channel)
}

func LoadChannel(config db.ChannelConfig) (err error) {
	var ch message.Channel
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
		err = ch.RefreshSession() // 保证session不空
		if err != nil {
			return
		}
	}

	channels[ch.Route] = ch
	return
}

func resolveChSecret(s string, category db.RouteCategory) (c message.ChannelCipher, err error) {
	var chSecret db.ChannelSecurity
	err = json.Unmarshal([]byte(s), &chSecret)
	if err != nil {
		return
	}

	var pubKey *rsa.PublicKey
	var priKey *rsa.PrivateKey
	if category == db.Subscriber {
		pubKey, err = cypher.LoadPublicKey([]byte(chSecret.RsaPem))
		if err != nil {
			return
		}
		c.RsaPubKey = pubKey // 通道中只加载公钥
	}
	if category == db.Publish {
		priKey, err = cypher.LoadPrivateKey([]byte(chSecret.RsaPem))
		if err != nil {
			return
		}
		c.RsaPriKey = priKey
	}

	return
}

// GetChannel 查询支持的通道
func GetChannel(name string) (ch message.Channel, err error) {
	ch, ok := channels[name]
	if !ok {
		err = errors.New("channel not found")
		return
	}
	return
}
