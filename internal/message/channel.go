package message

import (
	"crypto/rsa"
	"fmt"

	"gitee.com/lishimeng/event-bus/internal/db"
	"gitee.com/lishimeng/event-bus/internal/tls/cypher"
	"gitee.com/lishimeng/event-bus/internal/tls/session"
)

type Channel struct {
	Code     string           // 编号
	Name     string           // 名称
	Category db.RouteCategory // 路径类型
	Route    string           // 路由(目的地)不支持二级路由
	UseTls   bool             // 加密开关(不加密时message的payload与biz_payload相同)
	Cipher   ChannelCipher
	Callback string     // callback uri
	s        *session.S // publish通道复用session
}

func (ch *Channel) GetKey() string {
	return fmt.Sprintf("%d_%s", ch.Category, ch.Route)
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
