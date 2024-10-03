package models

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type eventHandleFunction func(client *WsClient)

type WsServer struct {
	server  websocket.Upgrader
	clients []*WsClient
	mu sync.Mutex
	handlers map[string]eventHandleFunction
}

type ServerResponse struct {
	Status string `json:"status"`
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
	ws.addClient(&wsClient)
	wsClient.Emit("reply", ServerResponse{Status: "Ok"})

	go ws.waitForMessage(&wsClient)
}

func (ws *WsServer) addClient(client *WsClient) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.clients = append(ws.clients, client)
}

func (ws *WsServer) removeClient(client *WsClient) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	for i, c := range ws.clients {
		if c.Id == client.Id {
			ws.clients = append(ws.clients[:i], ws.clients[i+1:]...)
			break
		}
	}
}

func (ws *WsServer) waitForMessage(client *WsClient) {
	conn := client.Conn

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Erro ao ler mensagem do cliente %d: %v", client.Id, err)
			ws.removeClient(client)
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

func (ws *WsServer) Emit(event string, message string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	for _, client := range ws.clients {
		if err := client.Emit(event, message); err != nil {
			log.Printf("Erro ao enviar mensagem para o cliente %d: %v", client.Id, err)
			client.Conn.Close()
			ws.removeClient(client)
		}
	}
}

func (ws *WsServer) setEventHandler(event string, handler eventHandleFunction) {
	ws.handlers[event] = handler
}

func (ws *WsServer) triggerEvent(event string, client *WsClient) {
	if handler, exists := ws.handlers[event]; exists {
		handler(client)
	}
}