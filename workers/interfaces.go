package workers

import "github.com/ThreeDotsLabs/watermill/message"

type MessageHandlerInterface interface {
	ReadMessages(messages <-chan *message.Message)
	Invoke(message *message.Message)
}
