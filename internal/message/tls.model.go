package message

import "crypto/rsa"

// ChannelCipher 通道的密钥
type ChannelCipher struct {
	RsaPubKey *rsa.PublicKey
	RsaPriKey *rsa.PrivateKey
}
