package models

import (
	"github.com/gorilla/websocket"
)

type WsClient struct {
	Conn *websocket.Conn
	Id uint
}

const lim = "|x1e"

func NewWsClient(connection *websocket.Conn, id uint) WsClient {
	return WsClient{
		Conn: connection,
		Id: id,
	}
}

func (client *WsClient) sendMessage(data []byte) error {
	return client.Conn.WriteMessage(websocket.TextMessage, data)
}

func (client *WsClient) Emit(event string, message string) error {
	combinedMessage := "4" + lim + event + lim + message
	return client.sendMessage([]byte(combinedMessage))
}