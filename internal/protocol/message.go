package protocol

import (
	"encoding/json"
	"fmt"
	"time"
)

type MessageType string

const (
	MessageTypeHello   MessageType = "HELLO"
	MessageTypeMessage MessageType = "MESSAGE"
	MessageTypeAck     MessageType = "ACK"
)

type Envelope struct {
	Type           MessageType `json:"type"`
	ID             string      `json:"id"`
	SenderID       string      `json:"sender_id"`
	Timestamp      time.Time   `json:"timestamp"`
	Payload        string      `json:"payload,omitempty"`
	AcknowledgesID string      `json:"acknowledges_id,omitempty"`
}

func EncodeMessage(message Envelope) ([]byte, error) {
	data, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("encode message error: %w", err)
	}
	return data, nil
}

func DecodeMessage(data []byte) (Envelope, error) {
	var message Envelope
	if err := json.Unmarshal(data, &message); err != nil {
		return Envelope{}, fmt.Errorf("decode message error: %w", err)
	}
	return message, nil
}
