package models

import (
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type     string            `json:"type"`
	Username string            `json:"user,omitempty"`
	User2    string            `json:"user2"`
	Content  string            `json:"content,omitempty"`
	Users    map[string]string `json:"users,omitempty"`
	Log      []string          `json:"log,omitempty"`
}

type Client struct {
	Conn     *websocket.Conn
	Username string
	LastPong time.Time
}
