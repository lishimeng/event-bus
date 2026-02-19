package db

import (
	"github.com/lishimeng/app-starter"
)

type RouteCategory int

type DataRecordState int

const (
	Subscriber RouteCategory = 1 // 自己的接收通道
	Publish    RouteCategory = 2 // 发布到第三方系统的通道
)

const (
	Init       DataRecordState = 1
	Processing DataRecordState = 2
	Success    DataRecordState = 3
	Failure    DataRecordState = 99
)

type DataRecord struct {
	app.Pk
	Code       string        `orm:"column(code);unique"`
	ReferCode  string        `orm:"column(refer_code)"`  // 关联的编号
	Source     string        `orm:"column(source)"`      // 端口(来源/目的地)
	Route      RouteCategory `orm:"column(route)"`       // 路由类型(到达/发出)
	Payload    string        `orm:"column(payload)"`     // 通信内容
	BizPayload string        `orm:"column(biz_payload)"` // 业务数据原文
	app.TableChangeInfo
}
