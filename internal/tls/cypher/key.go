package cypher

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

const (
	RsaKeySize = 2048
)

// GenerateKeys 密钥生成(固定2048)
func GenerateKeys() (priKey *rsa.PrivateKey, pubKey *rsa.PublicKey, err error) {
	priKey, err = rsa.GenerateKey(rand.Reader, RsaKeySize)
	if err != nil {
		return
	}
	pub := priKey.Public()
	pubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return
	}
	return
}

// SavePrivateKey 保存密钥(pkcs#8)
func SavePrivateKey(filename string, key *rsa.PrivateKey) (err error) {
	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return
	}
	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}
	return os.WriteFile(filename, pem.EncodeToMemory(privBlock), 0600)
}

// SavePublicKey 保存公钥(X.509)
func SavePublicKey(filename string, pubKey *rsa.PublicKey) error {
	pubASN1, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return err
	}
	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	}
	return os.WriteFile(filename, pem.EncodeToMemory(pubBlock), 0644)
}
func LoadPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return LoadPrivateKey(data)
}

func LoadPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
	}

	rsaPriv, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}
	return rsaPriv, nil
}

func LoadPublicKeyFromFile(filename string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return LoadPublicKey(data)
}

func LoadPublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}
