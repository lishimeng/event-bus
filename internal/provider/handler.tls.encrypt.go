package provider

import (
	"encoding/base64"
	"encoding/json"

	"gitee.com/lishimeng/event-bus/internal/channel"
	"gitee.com/lishimeng/event-bus/internal/db"
	"gitee.com/lishimeng/event-bus/internal/message"
	"gitee.com/lishimeng/event-bus/internal/tls/cypher"
	"gitee.com/lishimeng/event-bus/internal/tls/session"
)

// TlsEncryptHandler 加
var TlsEncryptHandler MessageHandler = func(msg *message.Message, ctx map[string]any) (err error) {
	ch, err := channel.GetChannel(msg.Route, db.PublishTo) // 用于发布的通道
	if err != nil {
		return
	}
	var nonce []byte
	var key []byte
	var aesKey []byte
	data := msg.Payload.Data
	dataBytes, err := base64.StdEncoding.DecodeString(data)
	if ch.UseTls {
		var s session.S
		nonce, err = base64.StdEncoding.DecodeString(msg.Payload.Nonce)
		if err != nil {
			return
		}
		key, err = base64.StdEncoding.DecodeString(msg.Payload.Key)
		if err != nil {
			return
		}
		aesKey, err = cypher.Decrypt(key, ch.Cipher.RsaPriKey)
		if err != nil {
			return
		}

		s, err = session.GenSession(aesKey, key)
		if err != nil {
			return
		}
		dataBytes, err = s.Decrypt(dataBytes, nonce)
		if err != nil {
			return
		}
	} else {

	}
	var biz message.BizMessage
	err = json.Unmarshal(dataBytes, &biz)
	if err != nil {
		return
	}
	ctx["biz"] = biz
	return
}
