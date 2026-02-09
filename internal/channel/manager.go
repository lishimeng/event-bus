package channel

import (
	"errors"

	"gitee.com/lishimeng/event-bus/internal/message"
)

// 通道格式为 route:channel(路由/通道实例)

// 订阅通道
var subscribers map[string]message.Channel

// 发布通道(一般只有一个)
var publishers map[string]message.Channel

func init() {
	subscribers = make(map[string]message.Channel)
	publishers = make(map[string]message.Channel)
}

func GetSubscriber(name string) (ch message.Channel, err error) {
	ch, ok := subscribers[name]
	if !ok {
		err = errors.New("channel not found")
		return
	}
	return
}

func GetPublisher(name string) (ch message.Channel, err error) {
	ch, ok := publishers[name]
	if !ok {
		err = errors.New("channel not found")
		return
	}
	return
}
