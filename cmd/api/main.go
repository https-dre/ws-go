package main

import (
	"log"
	"web-socket/pkg/models"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	wsServer := models.NewWsServer()

	router.GET("/ws", wsServer.HandleNewConnection)
	log.Println("SERVIDOR WEBSOCKET ESCUTANDO NA PORTA 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}