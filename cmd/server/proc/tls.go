package proc

import "crypto/rsa"

type ChannelCipher struct {
	RsaPubKey *rsa.PublicKey
	RsaPriKey *rsa.PrivateKey
}

var UserLocalCipher = false

// LocalCipher 本地密钥
var LocalCipher ChannelCipher
