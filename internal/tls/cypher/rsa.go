package cypher

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// Encrypt RSA加密（使用OAEP填充）
func Encrypt(plaintext []byte, pubKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, plaintext, nil)
}

// Decrypt RSA解密（使用OAEP填充）
func Decrypt(ciphertext []byte, privKey *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, ciphertext, nil)
}
