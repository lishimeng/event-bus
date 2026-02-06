package message

import "gitee.com/lishimeng/event-bus/internal/tls/session"

type Message struct {
	RequestId string          `json:"requestId,omitempty"`
	ReferId   string          `json:"referId,omitempty"`
	Route     string          `json:"route,omitempty"`
	Payload   session.Payload `json:"payload"`
}

type BizMessage struct {
	Action  string            `json:"action,omitempty"`
	Method  string            `json:"method,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Data    map[string]any    `json:"data,omitempty"`
}
