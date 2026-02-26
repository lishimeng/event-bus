package sdk

import (
	"encoding/base64"
	"encoding/json"
	"net/url"

	"gitee.com/lishimeng/event-bus/internal/message"
	"github.com/lishimeng/go-log"
)

const publishPath = "/api/v1/communication/publish" // 默认位置

// Request Payload字段必须是message.BizMessage序列化后的格式, 并经过base64转码
type Request struct {
	Payload string             `json:"payload,omitempty"` // 优先使用, 空白的时候从Biz里读数据
	Biz     message.BizMessage `json:"biz,omitempty"`     // payload有值的时候不生效
	Route   string             `json:"route,omitempty"`
	ReferId string             `json:"referId,omitempty"` // 作为主动回复时可标记原message_id
}

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type Client struct {
	Host string // server host
}

func (c *Client) genUrl(path string) string {
	dest, err := url.JoinPath(c.Host, path)
	if err != nil {
		return ""
	}
	return dest
}

func (c *Client) Publish(route string, msg message.BizMessage) (result Resp, err error) {
	method := msg.Method
	log.Info("method:%v", method)
	bs, err := json.Marshal(msg)
	if err != nil {
		return
	}
	payload := base64.StdEncoding.EncodeToString(bs)
	u := c.genUrl(publishPath)
	m := make(map[string]any)
	m["payload"] = payload
	m["route"] = route
	resp, err := GetDefault().Post(u, m)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return
	}
	return
}
