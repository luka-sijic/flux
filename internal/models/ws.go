package models

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Send     chan []byte
	Username string
	LastPong time.Time
}
