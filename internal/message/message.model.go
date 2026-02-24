package message

import "gitee.com/lishimeng/event-bus/internal/tls/session"

type Message struct {
	RequestId string          `json:"requestId,omitempty"`
	ReferId   string          `json:"referId,omitempty"`
	Route     string          `json:"route,omitempty"`
	Payload   session.Payload `json:"payload"`
	Biz       BizMessage      `json:"-"` // 业务数据,不参与序列化
}

func (m *Message) Decrypt(ch Channel) (biz BizMessage, err error) {
	if ch.UseTls {

	}
	return
}

// BizMessage 业务消息体
type BizMessage struct {
	Action  string            `json:"action,omitempty"`
	Method  string            `json:"method,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Data    map[string]any    `json:"data,omitempty"`
}
