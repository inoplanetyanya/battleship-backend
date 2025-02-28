package app

import (
	wsserver "battleship/pkg/websocket"
	"net/http"

	"golang.org/x/net/websocket"
)

func Run() {
	server := wsserver.NewServer(1024)
	http.Handle("/ws", websocket.Handler(server.HandleWS))
	http.ListenAndServe(":8080", nil)
}
