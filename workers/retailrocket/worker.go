package retailrocket

import (
	"github.com/ThreeDotsLabs/watermill/message"
)

type RetailRocket struct {
	HttpClient *Client
}

func (receiver *RetailRocket) ReadMessages(messages <-chan *message.Message) {
	for msg := range messages {
		receiver.Invoke(msg)
	}
}

func (receiver *RetailRocket) Invoke(msg *message.Message) {
	postMessage, err := NewMessage(msg)
	if err != nil {
		return
	}

	if receiver.HttpClient.Post(postMessage) != nil {
		msg.Nack()
		return
	}

	msg.Ack()
}
