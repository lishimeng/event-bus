package provider

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"gitee.com/lishimeng/event-bus/internal/channel"
	"gitee.com/lishimeng/event-bus/internal/message"
	"gitee.com/lishimeng/event-bus/internal/tls/cypher"
	"gitee.com/lishimeng/event-bus/internal/tls/session"
)

type MessageHandler func(msg message.Message, ctx map[string]any) (err error)

const (
	MaxMessageBuffer = 32
)

var ChannelChkHandler MessageHandler = func(msg message.Message, ctx map[string]any) (err error) {
	_, err = channel.GetChannel(msg.Route)
	if err != nil {
		return
	}
	return
}

var queue = make(chan message.Message, MaxMessageBuffer)

// DataRecordMsgHandler 存储
var DataRecordMsgHandler MessageHandler = func(msg message.Message, ctx map[string]any) (err error) {
	select {
	case queue <- msg:
	case <-time.After(time.Millisecond * 10):
		return
	}
	return
}

var TlsDecryptHandler MessageHandler = func(msg message.Message, ctx map[string]any) (err error) {
	ch, err := channel.GetChannel(msg.Route)
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
