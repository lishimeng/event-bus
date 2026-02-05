package proc

import (
	"encoding/json"

	"gitee.com/lishimeng/event-bus/internal/tls/session"
	"github.com/lishimeng/go-log"
)

type MessageCreateOpt struct {
	Id       string
	ParentId string
}

type MessageCreateFunc func(opt *MessageCreateOpt)

var WithId = func(id string) MessageCreateFunc {
	return func(opt *MessageCreateOpt) {
		opt.Id = id
	}
}

var WithParentId = func(id string) MessageCreateFunc {
	return func(opt *MessageCreateOpt) {
		opt.ParentId = id
	}
}

type Message struct {
	RequestId string          `json:"requestId,omitempty"`
	ReferId   string          `json:"referId,omitempty"`
	Route     string          `json:"route,omitempty"`
	Payload   session.Payload `json:"payload"`
}

func Create(destination string, payload any, opts ...MessageCreateFunc) (m Message, err error) {
	// 消息创建业务
	var p session.Payload
	var opt MessageCreateOpt
	for _, o := range opts {
		o(&opt)
	}
	m.RequestId = opt.Id
	m.ReferId = opt.ParentId
	m.Route = destination
	if len(m.RequestId) == 0 {
		// TODO gen id
	}

	var nonce []byte
	var ciphertext []byte

	bs, err := json.Marshal(payload)
	if err != nil {
		return
	}
	// get destination
	ch, err := GetChannel(destination)
	if err != nil {
		log.Info("not found channel[%s]", destination)
		return
	}

	if ch.UseTls {
		var s = ch.GetSession()
		log.Info("use tls")
		ciphertext, nonce, err = s.Encrypt(bs)
		if err != nil {
			log.Info("encrypt fail")
			log.Info(err)
			return
		}
		p = s.GenData(s.AesKey, nonce, ciphertext)
	} else {
		p.Data = string(bs)
	}

	m.Payload = p
	return
}
