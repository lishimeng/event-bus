package message

import "crypto/rsa"

type ChannelCipher struct {
	RsaPubKey *rsa.PublicKey
	RsaPriKey *rsa.PrivateKey
}
