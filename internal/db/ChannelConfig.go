package db

import (
	"encoding/base64"
	"encoding/json"

	"github.com/lishimeng/app-starter"
)

// ChannelConfig 通道配置,
type ChannelConfig struct {
	app.Pk
	Code        string        `gorm:"column:code;uniqueIndex"`
	Name        string        `gorm:"column:name"`
	Category    RouteCategory `gorm:"column:category"` // 路由方向
	Router      string        `gorm:"column:router"`   // 路由路径
	UseSecurity int           `gorm:"column:use_security"`
	Security    string        `gorm:"column:security"` // 公钥
	Callback    string        `gorm:"column:callback"` // 回调配置
	app.TableChangeInfo
}

// ChannelSecurity 通道密钥, 订阅通道只需要公钥, 发布通道只需要私钥
type ChannelSecurity struct {
	RsaKey string `json:"rsaKey,omitempty"`
	RsaPem string `json:"rsaPem,omitempty"`
}

func (cs *ChannelSecurity) Marshal() string {
	bs, err := json.Marshal(cs)
	if err != nil {
		return ""
	}
	s := base64.StdEncoding.EncodeToString(bs)
	return s
}

func (cs *ChannelSecurity) Unmarshal(s string) (err error) {
	bs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, &cs)
	if err != nil {
		return
	}
	return
}
