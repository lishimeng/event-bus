package proc

import (
	"gitee.com/lishimeng/event-bus/internal/channel"
	"gitee.com/lishimeng/event-bus/internal/message"
	"gitee.com/lishimeng/event-bus/internal/tls/cypher"
	"gitee.com/lishimeng/event-bus/internal/tls/session"
	"github.com/lishimeng/go-log"
)

func Callback(m message.Message) (err error) {
	// 回调业务
	log.Info("callback message")
	route := m.Route
	ch, err := channel.GetPublisher(route)
	if err != nil {
		log.Info(err)
		return
	}
	var biz = m.Biz
	log.Info("callback_uri: %s", ch.Callback)
	log.Info("[%s]%s", biz.Method, biz.Action)
	for key, value := range biz.Data {
		log.Info("%s=%v", key, value)
	}
	return
}

func decodeData(key, nonce, data []byte) (plain []byte, err error) {
	var s session.S
	var aesKey []byte
	aesKey, err = cypher.Decrypt(key, LocalCipher.RsaPriKey)
	if err != nil {
		return
	}
	s, err = session.GenSession(aesKey, key)
	if err != nil {
		return
	}
	plain, err = s.Decrypt(data, nonce)
	return
}
