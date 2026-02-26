package message

import "gitee.com/lishimeng/event-bus/internal/tls/session"

type Message struct {
	RequestId string          `json:"requestId,omitempty"`
	ReferId   string          `json:"referId,omitempty"`
	Route     string          `json:"route,omitempty"`  // 发送的目的地
	Source    string          `json:"source,omitempty"` // 来源地(可选)
	Payload   session.Payload `json:"payload"`
	Biz       BizMessage      `json:"-"` // 业务数据,不参与序列化
}

func (m *Message) Decrypt(ch Channel) (biz BizMessage, err error) {
	if ch.UseTls {

	}
	return
}

// BizCallback 自动回复消息的配置
//
// 如果配置了callback, 会将执行后的response放在data中创建一个biz_message发送到原位置
//
// message中需要有来源地[source]
type BizCallback struct {
	CallbackAction string `json:"callbackAction,omitempty"` //
	CallbackMethod string `json:"callbackMethod,omitempty"`
}

// BizMessage 业务消息体
type BizMessage struct {
	BizCallback
	Action  string            `json:"action,omitempty"`
	Method  string            `json:"method,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Data    map[string]any    `json:"data,omitempty"`
}
