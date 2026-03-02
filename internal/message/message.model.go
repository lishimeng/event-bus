package message

import (
	"github.com/lishimeng/event-bus/internal/tls/session"
	"github.com/lishimeng/event-bus/sdk"
)

type Message struct {
	RequestId string          `json:"requestId,omitempty"`
	ReferId   string          `json:"referId,omitempty"`
	Route     string          `json:"route,omitempty"`  // 发送的目的地
	Source    string          `json:"source,omitempty"` // 来源地(可选)
	Payload   session.Payload `json:"payload"`
	Biz       sdk.BizMessage  `json:"-"` // 业务数据,不参与序列化
}

func (m *Message) Decrypt(ch Channel) (biz sdk.BizMessage, err error) {
	if ch.UseTls {

	}
	return
}
