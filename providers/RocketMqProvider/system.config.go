package RocketMqProvider

import (
	"errors"
	"fmt"
)

type RocketPublisher struct {
	MessageGroup string   `json:"messageGroup,omitempty"`
	Topics       []string `json:"topics,omitempty"`
}

type RocketSubscriber struct {
	ConsumerGroup string `json:"consumerGroup,omitempty"`
	Topic         string `json:"topic,omitempty"`
}

type RmqConfig struct {
	Endpoint    string             `json:"endpoint,omitempty"`
	AppId       string             `json:"appId,omitempty"`
	Secret      string             `json:"secret,omitempty"`
	Publisher   RocketPublisher    `json:"publisher,omitempty"`
	Subscribers []RocketSubscriber `json:"subscribers,omitempty"`
}

func (rc *RmqConfig) GetSubscriber(topic string) (s RocketSubscriber, err error) {
	for _, sub := range rc.Subscribers {
		if sub.Topic == topic {
			s = sub
			return
		}
	}
	err = errors.New(fmt.Sprintf("subscriber %s not exists", topic))
	return
}
