package message

import (
	"encoding/base64"
	"encoding/json"

	"github.com/lishimeng/event-bus/internal/tls/cypher"
	"github.com/lishimeng/event-bus/internal/tls/session"
	"github.com/lishimeng/event-bus/sdk"
	"github.com/lishimeng/go-log"
)

func Encrypt(biz sdk.BizMessage, target Channel) (payload session.Payload, err error) {
	bs, err := json.Marshal(biz)
	if err != nil {
		return
	}
	body := base64.StdEncoding.EncodeToString(bs)
	if !target.UseTls {
		payload.Data = body // 明文, 只有data
		return
	}
	s := target.GetSession()
	ciphertext, nonce, err := s.Encrypt(bs)
	if err != nil {
		return
	}
	p := s.GenData(s.AesKey, nonce, ciphertext)
	payload.Key = p.Key
	payload.Data = p.Data
	payload.NonceLen = p.NonceLen
	payload.Nonce = p.Nonce
	payload.Padding = p.Padding
	payload.TagLen = p.TagLen
	return
}

func Decrypt(payload session.Payload, receiver Channel) (biz sdk.BizMessage, err error) {
	body := payload.Data
	dataBytes, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return
	}
	if receiver.UseTls { // 需要解密
		dataBytes, err = decode(payload, receiver)
		if err != nil {
			log.Info("decode biz body fail")
			log.Info(err)
			return
		}
	}
	err = json.Unmarshal(dataBytes, &biz)
	if err != nil {
		log.Info("unmarshal biz body fail")
		log.Info(err)
		return
	}
	return
}

// 解密
func decode(payload session.Payload, receiver Channel) (dataBytes []byte, err error) {
	var s session.S
	var nonce []byte
	var key []byte
	var aesKey []byte
	data := payload.Data
	dataBytes, err = base64.StdEncoding.DecodeString(data)

	nonce, err = base64.StdEncoding.DecodeString(payload.Nonce)
	if err != nil {
		return
	}
	key, err = base64.StdEncoding.DecodeString(payload.Key)
	if err != nil {
		return
	}
	aesKey, err = cypher.Decrypt(key, receiver.Cipher.RsaPriKey)
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

	return
}
