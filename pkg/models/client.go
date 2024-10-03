package models

import (
	"encoding/json"
	"errors"

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

func (client *WsClient) Emit(event string, data interface{}) error {
	var message []byte
	var messageType string
	var err error

	switch v := data.(type) {
	case string:
		message = []byte(v)
		messageType = "string"
	default:
		messageType = "json"
		message, err = json.Marshal(data)
		if err != nil {
			return errors.New("erro ao serializar o dado para JSON")
		}
	}
	combinedMessage := []byte("4" + lim + event + lim + messageType + lim + "\n")
	combinedMessage = append(combinedMessage, message...)

	return client.sendMessage(combinedMessage)
}