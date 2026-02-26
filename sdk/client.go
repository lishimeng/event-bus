package sdk

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

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

func (r *Request) WithReferId(id string) *Request {
	r.ReferId = id
	return r
}

func (r *Request) WithPayload(payload string) *Request {
	r.Payload = payload
	return r
}

func (r *Request) WithBiz(biz message.BizMessage) *Request {
	r.Biz = biz
	return r
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

func (c *Client) CreateRequest(route string) (r Request) {
	r.Route = route
	return
}

func (c *Client) Publish(msg Request) (result Resp, err error) {
	method := msg.Biz.Method
	log.Info("method:%v", method)
	method = strings.ToUpper(method)
	switch method {
	case "POST":
	case "GET":
	case "PUT":
	case "DELETE":
	default:
		err = errors.New("invalid method")
		return
	}
	u := c.genUrl(publishPath)
	resp, err := GetDefault().Post(u, msg)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return
	}
	return
}
