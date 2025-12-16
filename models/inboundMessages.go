package models

import "encoding/json"

type BaseInbound struct {
	MsgType string       `json:"type"`
	Data json.RawMessage `json:"data"`
}

type LoginInbound struct {
	Token string `json:"token"`
}