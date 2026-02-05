package proc

import (
	"encoding/base64"
	"encoding/json"

	"gitee.com/lishimeng/event-bus/internal/tls/cypher"
	"gitee.com/lishimeng/event-bus/internal/tls/session"
	"github.com/lishimeng/go-log"
)

func Callback(m Message) (err error) {
	// 回调业务
	log.Info("callback message")
	bs, _ := json.Marshal(m)
	log.Info(string(bs))

	source := m.Route
	key, err := base64.StdEncoding.DecodeString(m.Payload.Key)
	nonce, err := base64.StdEncoding.DecodeString(m.Payload.Nonce)
	data, err := base64.StdEncoding.DecodeString(m.Payload.Data)

	var payload []byte
	if UserLocalCipher {
		payload, err = decodeData(key, nonce, data)
		if err != nil {
			log.Info("decode data failed")
			return
		}
	} else {
		payload = data
	}

	log.Info("<-%s [%s]", source, string(payload))
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
