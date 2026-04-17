package RocketMqProvider

import "strings"

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

func (rc *RmqConfig) PublisherMessageGroup() string {
	if strings.TrimSpace(rc.Publisher.MessageGroup) != "" {
		return strings.TrimSpace(rc.Publisher.MessageGroup)
	}
	return "event-bus-publisher"
}

func (rc *RmqConfig) ConsumerGroup() string {
	for _, sub := range rc.Subscribers {
		if strings.TrimSpace(sub.ConsumerGroup) != "" {
			return strings.TrimSpace(sub.ConsumerGroup)
		}
	}
	return ""
}

func (rc *RmqConfig) GetSubscriber(topic string) (s RocketSubscriber, ok bool) {
	topic = strings.TrimSpace(topic)
	for _, sub := range rc.Subscribers {
		if strings.TrimSpace(sub.Topic) == topic {
			return sub, true
		}
	}
	if len(rc.Subscribers) > 0 {
		return rc.Subscribers[0], true
	}
	return RocketSubscriber{}, false
}
