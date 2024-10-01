package models
import (
	"log"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsServer struct {
	server  websocket.Upgrader
	clients []*WsClient
}

func NewWsServer() *WsServer {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &WsServer{
		server:  upgrader,
		clients: make([]*WsClient, 0),
	}
}

func (ws *WsServer) HandleNewConnection(c *gin.Context) {
	conn, err := ws.server.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Println(err)
		c.String(400, err.Error())
		return
	}

	wsClient := NewWsClient(conn, uint(len(ws.clients)+1))
	ws.AddClient(&wsClient)

	go ws.WaitForMessage(&wsClient)

}

func (ws *WsServer) AddClient(client *WsClient) {
	ws.clients = append(ws.clients, client)
}

func (ws *WsServer) RemoveClient(client *WsClient) {
	for i, c := range ws.clients {
		if c.Id == client.Id {
			ws.clients = append(ws.clients[:i], ws.clients[i+1:]...)
			break
		}
	}
}

func (ws *WsServer) WaitForMessage(client *WsClient) {
	conn := client.Conn

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Erro ao ler mensagem do cliente %d: %v", client.Id, err)
			ws.RemoveClient(client)
			conn.Close()
			break
		}

		switch msgType {
		case websocket.TextMessage:
			log.Printf("-> %d: %s", client.Id, string(msg))
		case websocket.BinaryMessage:
			log.Printf("b-> %d", client.Id)
		default:
			log.Printf("unknown -> %d", client.Id)
		}
	}
}

func (ws *WsServer) SendMessage(client *WsClient, data []byte) error {
	err := client.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}
	return nil
}
