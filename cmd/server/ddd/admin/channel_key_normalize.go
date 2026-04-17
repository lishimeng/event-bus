package admin

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
)

// stripBase64Noise 去掉 JSON/粘贴带来的换行与空白，便于 StdEncoding 解码。
func stripBase64Noise(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '\n', '\r', ' ', '\t':
			continue
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// normalizePrivateKeyToPEM 支持：PEM 原文、或 PKCS#8/PKCS#1 的 DER 二进制。
// 统一为 RSA PRIVATE KEY，确保与后端加载逻辑一致。
func normalizePrivateKeyToPEM(raw []byte) ([]byte, error) {
	t := bytes.TrimSpace(raw)
	if len(t) == 0 {
		return nil, fmt.Errorf("empty private key")
	}
	if bytes.HasPrefix(t, []byte("-----BEGIN")) {
		block, _ := pem.Decode(t)
		if block == nil || block.Type != "RSA PRIVATE KEY" {
			return nil, fmt.Errorf("private key PEM type must be RSA PRIVATE KEY")
		}
		return t, nil
	}
	if k, err := x509.ParsePKCS8PrivateKey(t); err == nil {
		rsaK, ok := k.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not RSA")
		}
		der, err := x509.MarshalPKCS8PrivateKey(rsaK)
		if err != nil {
			return nil, err
		}
		return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}), nil
	}
	if rsaK, err := x509.ParsePKCS1PrivateKey(t); err == nil {
		der := x509.MarshalPKCS1PrivateKey(rsaK)
		return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}), nil
	}
	return nil, fmt.Errorf("unsupported private key: need PEM or RSA PKCS#8/PKCS#1 DER")
}

// normalizePublicKeyToPEM 支持：PEM 原文、或 PKIX 公钥 DER。
func normalizePublicKeyToPEM(raw []byte) ([]byte, error) {
	t := bytes.TrimSpace(raw)
	if len(t) == 0 {
		return nil, fmt.Errorf("empty public key")
	}
	if bytes.HasPrefix(t, []byte("-----BEGIN")) {
		return t, nil
	}
	pub, err := x509.ParsePKIXPublicKey(t)
	if err != nil {
		return nil, fmt.Errorf("public key: %w", err)
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not RSA")
	}
	der, err := x509.MarshalPKIXPublicKey(rsaPub)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}), nil
}
