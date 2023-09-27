package model

import (
	"encoding/json"
	"sync/atomic"
)

type Message struct {
	Id        string `json:"id,omitempty"`
	Content   string `json:"content,omitempty"`
	Channel   string `json:"channel,omitempty"`
	Command   string `json:"command,omitempty"`
	CreatedAt string `json:"created_at"`
	Sender    string `json:"sender" bson:"sender,omitempty"`
	Err       string `json:"err,omitempty"`
	json      atomic.Value
}

// Hash returns the transaction hash.
func (msg *Message) JSONEncode() []byte {
	if json := msg.json.Load(); json != nil {
		return json.([]byte)
	}
	jsonBytes, _ := json.Marshal(msg)
	msg.json.Store(jsonBytes)
	return jsonBytes
}
