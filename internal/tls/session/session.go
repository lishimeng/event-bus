package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

const (
	NonceSile  = 12 // byte
	AesKeySize = 32
)

type S struct {
	aesKey   []byte
	AesKey   []byte
	TagLen   int
	Padding  string
	NonceLen int
	gcm      *cipher.AEAD
}

// GenData 输出格式化数据
func (s *S) GenData(key []byte, nonce []byte, ciphertext []byte) (p Payload) {
	p.Key = base64.StdEncoding.EncodeToString(key)
	p.Nonce = base64.StdEncoding.EncodeToString(nonce)
	p.Data = base64.StdEncoding.EncodeToString(ciphertext)
	p.NonceLen = s.NonceLen
	p.TagLen = s.TagLen
	p.Padding = s.Padding
	return
}

// GenNonce 随机数产生器
func (s *S) genNonce() ([]byte, error) {
	var err error
	nonce := make([]byte, NonceSile)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

// Encrypt 加密
func (s *S) Encrypt(plaintext []byte) (ciphertext []byte, nonce []byte, err error) {

	nonce, err = s.genNonce()
	if err != nil {
		return
	}
	var gcm = *s.gcm
	ciphertext = gcm.Seal(nil, nonce, plaintext, nil)
	return
}

// Decrypt 解密
func (s *S) Decrypt(ciphertext []byte, nonce []byte) ([]byte, error) {
	var gcm = *s.gcm
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext, err
}

// GenAesKey 产生aes_key(随机模式)
func GenAesKey() ([]byte, error) {
	var err error
	key := make([]byte, AesKeySize)
	if _, err = io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

func GenSession(key []byte, encodedKey []byte) (s S, err error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	s.aesKey = key
	s.AesKey = encodedKey
	s.TagLen = 16 * 8
	s.NonceLen = NonceSile * 8
	s.Padding = "NoPadding"
	s.gcm = &gcm
	return
}
