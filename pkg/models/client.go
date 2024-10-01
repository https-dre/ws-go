package models

import "github.com/gorilla/websocket"

type WsClient struct {
	Conn *websocket.Conn
	Id uint
}

func NewWsClient(connection *websocket.Conn, id uint) WsClient {
	return WsClient{
		Conn: connection,
		Id: id,
	}
}