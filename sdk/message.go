package sdk

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
	Action  string            `json:"apiPath,omitempty"`
	Method  string            `json:"method,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Data    map[string]any    `json:"params,omitempty"`
}
