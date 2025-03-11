package handler

import (
	"battleship/pkg/common"
	"battleship/pkg/service"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
)

type WebSocketHandler struct {
	service *service.Service
}

func NewWebSocketHandler(service *service.Service) *WebSocketHandler {
	return &WebSocketHandler{service: service}
}

func (h *WebSocketHandler) InitRoutes(router *http.ServeMux) {
	router.HandleFunc("/ws", websocket.Handler(h.hws).ServeHTTP)
}

func (h *WebSocketHandler) hws(ws *websocket.Conn) {
	log.Println("[ws] New client connection: ", ws.RemoteAddr())
	h.readLoop(ws)

}

func (h *WebSocketHandler) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("[ws] Connection closed", ws.RemoteAddr())
				break
			}
			log.Fatal("[ws] Read error: ", err)
			continue
		}
		msg := string(buf[:n])

		sm := strings.Split(msg, " ")
		if len(sm) < 2 {
			log.Println("[ws] Invalid message format:", msg)
			continue
		}

		if sm[0] == "/chat" {
			gr := h.service.Game.GameRoomList().MapByConn[ws]
			if gr == nil {
				continue
			}

			gr.Chat(ws, strings.Join(sm[1:], " "))
			continue
		}

		userId, err := h.service.Authorization.ParseToken(sm[1])
		if err != nil {
			log.Println("[ws] Token parse error:", err)
			continue
		}
		log.Printf("[ws] message from user with id %d: %s\n", userId, msg)

		if sm[0] == "/connect" {
			user, err := h.service.Authorization.GetUserByToken(sm[1])
			if err != nil {
				log.Println("[ws] GetUserByToken error:", err)
				continue
			}

			p := common.Player{
				User: user,
				Conn: ws,
			}

			g := h.service.Game

			if err := g.AddPlayerToQueue(p); err != nil {
				log.Println("[ws] AddPlayerToQueue error:", err)
			}
		}

		if sm[0] == "/disconnect" {
			user, err := h.service.Authorization.GetUserByToken(sm[1])
			if err != nil {
				log.Println("[ws] GetUserByToken error:", err)
				continue
			}

			p := common.Player{
				User: user,
				Conn: ws,
			}

			g := h.service.Game

			rp, err := g.RemovePlayerFromQueue(p)
			if err != nil {
				log.Println("[ws] RemovePlayerFromQueue error:", err)
				continue
			} else {
				log.Printf("[ws] player removed from queue %v\n", rp)
			}
		}
	}
}
