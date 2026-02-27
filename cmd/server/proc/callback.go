package proc

import (
	"errors"
	"net/url"
	"strings"

	"gitee.com/lishimeng/event-bus/internal/channel"
	"gitee.com/lishimeng/event-bus/internal/message"
	"gitee.com/lishimeng/event-bus/internal/tls/cypher"
	"gitee.com/lishimeng/event-bus/internal/tls/session"
	"gitee.com/lishimeng/event-bus/sdk"
	"github.com/lishimeng/go-log"
)

func Callback(m message.Message) (err error) {
	// 回调业务
	log.Info("callback message")
	route := m.Route
	ch, err := channel.GetChannel(route)
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
	var action string
	action, err = url.JoinPath(ch.Callback, biz.Action)
	if err != nil {
		log.Info(err)
		return
	}
	var method = biz.Method
	var resp sdk.Response
	method = strings.ToUpper(method)
	switch method {
	case "POST":
		resp, err = sdk.GetDefault().Post(action, biz.Data)
	case "PUT":
	case "DELETE":
	case "GET":
		resp, err = sdk.GetDefault().Get(action)
	default:
		err = errors.New("invalid method")
		return
	}
	if err != nil {
		log.Info(err)
		return
	}
	log.Info("callback response: %d", resp.StatusCode)
	log.Info("callback response: %s", string(resp.Body))

	if len(m.Source) > 0 && len(biz.BizCallback.CallbackAction) > 0 {
		// TODO 产生callback message
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
