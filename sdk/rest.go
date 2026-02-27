package sdk

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// RESTClient 封装 HTTP 客户端
type RESTClient struct {
	client *http.Client
}

type Response struct {
	StatusCode int
	Body       []byte
}

// NewRESTClient 创建一个新的 REST 客户端
func NewRESTClient(timeout time.Duration) *RESTClient {
	return &RESTClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

var defaultRestClient *RESTClient

func init() {
	defaultRestClient = NewRESTClient(time.Second * 15)
}

func GetDefault() *RESTClient {
	return defaultRestClient
}

// Get 发起 GET 请求
func (c *RESTClient) Get(url string) (r Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	r.StatusCode = resp.StatusCode
	r.Body = body
	return
}

// Post 发起 POST 请求，body 为任意可序列化为 JSON 的对象
func (c *RESTClient) Post(url string, body interface{}) (r Response, err error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	r.StatusCode = resp.StatusCode
	r.Body = respBody
	return
}

// Put 发起 PUT 请求，body 为任意可序列化为 JSON 的对象
func (c *RESTClient) Put(url string, body any) (r Response, err error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	r.StatusCode = resp.StatusCode
	r.Body = respBody
	return
}
