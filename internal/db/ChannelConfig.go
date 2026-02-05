package db

import (
	"github.com/lishimeng/app-starter"
)

// ChannelConfig 通道配置,
type ChannelConfig struct {
	app.Pk
	Code        string `orm:"column(code);unique"`
	Name        string `orm:"column(name)"`
	Router      string `orm:"column(router)"`
	UseSecurity int    `orm:"column(use_security)"`
	Security    string `orm:"column(security)"` // 公钥
	app.TableChangeInfo
}

type ChannelSecurity struct {
	RsaKey string `json:"rsaKey,omitempty"`
	RsaPem string `json:"rsaPem,omitempty"`
}
