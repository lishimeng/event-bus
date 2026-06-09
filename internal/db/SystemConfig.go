package db

import "github.com/lishimeng/app-starter"

const (
	SysLocalSecret = "local_secret"
	SysRmqConfig   = "rmq_config"
)

type SysConfig struct {
	app.Pk
	Name   string `gorm:"column:name"`   // 配置名称
	Config string `gorm:"column:config"` // base64格式配置内容
	app.TableChangeInfo
}

// LocalSecurity 本地密钥
type LocalSecurity struct {
	RsaKey string `json:"rsaKey,omitempty"`
	RsaPem string `json:"rsaPem,omitempty"`
}
