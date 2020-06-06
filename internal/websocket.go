package internal

import (
	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	Clients map[*websocket.Conn]string
}

func NewWebsocketClient() *WebsocketClient {
	client := &WebsocketClient{}
	client.Clients = make(map[*websocket.Conn]string)
	return client
}
