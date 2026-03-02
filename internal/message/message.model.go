package message

import (
	"github.com/lishimeng/event-bus/internal/tls/session"
	"github.com/lishimeng/event-bus/sdk"
)

type Message struct {
	RequestId string          `json:"messageId,omitempty"`
	ReferId   string          `json:"referenceId,omitempty"`
	Route     string          `json:"topic,omitempty"`  // 发送的目的地
	Source    string          `json:"source,omitempty"` // 来源地(可选)
	Payload   session.Payload `json:"payload"`
	Biz       sdk.BizMessage  `json:"-"` // 业务数据,不参与序列化
}
