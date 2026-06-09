package db

import (
	"github.com/lishimeng/app-starter"
)

type RouteCategory int

type DataRecordState int

const (
	Subscribe RouteCategory = 1 // 自己的接收通道
	PublishTo RouteCategory = 2 // 发布到第三方系统的通道
)

var categories = make(map[RouteCategory]string)

func init() {
	categories[Subscribe] = "Subscribe"
	categories[PublishTo] = "PublishTo"
}

func (rc RouteCategory) String() string {

	name, ok := categories[rc]
	if !ok {
		name = "unknown"
	}
	return name
}

const (
	Init       DataRecordState = 1
	Processing DataRecordState = 2
	Success    DataRecordState = 3
	Failure    DataRecordState = 99
)

type DataRecord struct {
	app.Pk
	Code       string        `gorm:"column:code;uniqueIndex"`
	ReferCode  string        `gorm:"column:refer_code"`  // 关联的编号
	Source     string        `gorm:"column:source"`      // 端口(来源/目的地)
	Route      RouteCategory `gorm:"column:route"`       // 路由类型(到达/发出)
	Payload    string        `gorm:"column:payload"`     // 通信内容
	BizPayload string        `gorm:"column:biz_payload"` // 业务数据原文
	app.TableChangeInfo
}
