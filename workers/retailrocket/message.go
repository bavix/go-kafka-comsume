package retailrocket

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Message struct {
	Method string `json:"method"`
	Body   string `json:"body"`
}

func NewMessage(message *message.Message) (*Message, error) {
	postMessage := &Message{}
	err := json.Unmarshal(message.Payload, postMessage)
	if err != nil {
		return nil, err
	}

	return postMessage, nil
}
