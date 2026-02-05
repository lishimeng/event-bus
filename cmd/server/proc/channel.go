package proc

import (
	"crypto/rsa"
	"encoding/json"
	"errors"

	"gitee.com/lishimeng/event-bus/internal/db"
	"gitee.com/lishimeng/event-bus/internal/tls/cypher"
	"gitee.com/lishimeng/event-bus/internal/tls/session"
)

type Channel struct {
	Code   string // 编号
	Name   string // 名称
	Route  string // 路由(目的地)不支持二级路由
	UseTls bool   // 加密开关(不加密时message的payload与biz_payload相同)
	Cipher ChannelCipher
	s      *session.S
}

func (ch *Channel) GetSession() *session.S {
	return ch.s
}

func (ch *Channel) RefreshSession() (err error) {

	s, err := createSession(ch.Cipher.RsaPubKey)
	if err != nil {
		return
	}
	ch.s = s
	return
}

var channels map[string]Channel // route:channel

func init() {
	channels = make(map[string]Channel)
}

func LoadChannel(config db.ChannelConfig) (err error) {
	var ch Channel
	ch.Code = config.Code
	ch.Name = config.Name
	ch.Route = config.Router
	if config.UseSecurity == 1 {
		ch.UseTls = true
		ch.Cipher, err = resolveChSecret(config.Security)
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

func createSession(pubKey *rsa.PublicKey) (s *session.S, err error) {
	aesKey, err := session.GenAesKey()
	if err != nil {
		return
	}
	cipheredAesKey, err := cypher.Encrypt(aesKey, pubKey)
	if err != nil {
		return
	}
	sessionIns, err := session.GenSession(aesKey, cipheredAesKey)
	if err != nil {
		return
	}
	s = &sessionIns
	return
}

func resolveChSecret(s string) (c ChannelCipher, err error) {
	var chSecret db.ChannelSecurity
	err = json.Unmarshal([]byte(s), &chSecret)
	if err != nil {
		return
	}

	pubKey, err := cypher.LoadPublicKey([]byte(chSecret.RsaPem))
	if err != nil {
		return
	}
	c.RsaPubKey = pubKey // 通道中只加载公钥
	return
}

// GetChannel 查询支持的通道
func GetChannel(name string) (ch Channel, err error) {
	ch, ok := channels[name]
	if !ok {
		err = errors.New("channel not found")
		return
	}
	return
}
