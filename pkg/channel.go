package tc

import (
	"github.com/gorilla/websocket"
)

type Channel struct {
	channel string
	ws      *websocket.Conn
}

func NewChannel(channel string, ws *websocket.Conn) *Channel {
	return &Channel{
		channel,
		ws,
	}
}
