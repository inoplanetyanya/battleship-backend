package common

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"golang.org/x/net/websocket"
)

type GameRoom struct {
	Players   map[*websocket.Conn]User
	Mu        sync.Mutex
	RoomID    string
	BoardSize int
}

func NewGameRoom() *GameRoom {
	return &GameRoom{
		Players:   make(map[*websocket.Conn]User),
		Mu:        sync.Mutex{},
		RoomID:    "room1",
		BoardSize: 10,
	}
}

func (g *GameRoom) AddPlayer(user User, conn *websocket.Conn) {
	g.Mu.Lock()
	defer g.Mu.Unlock()
	g.Players[conn] = user

	for _, p := range g.Players {
		if p == g.Players[conn] {
			fmt.Println("Player already in room")
		}
		return
	}

	fmt.Println("Player added")
}

func (g *GameRoom) RemovePlayer(conn *websocket.Conn) {
	g.Mu.Lock()
	delete(g.Players, conn)
	defer g.Mu.Unlock()
	fmt.Println("Player removed")
}

func (g *GameRoom) StartGame() {
	for player := range g.Players {
		player.Write([]byte("game start"))
	}
}

type ChatMessage struct {
	From    string `json:"from"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (g *GameRoom) Chat(from *websocket.Conn, msg string) {
	fromUser := g.Players[from]

	chatMsg := ChatMessage{
		From:    fromUser.Username,
		Message: string(msg),
		Type:    "chat",
	}

	jsonMsg, err := json.Marshal(chatMsg)
	if err != nil {
		log.Printf("Failed to marshal chat message: %v", err)
		return
	}

	for conn := range g.Players {
		if _, err := conn.Write(jsonMsg); err != nil {
			log.Printf("Failed to send message to player: %v", err)
		}
	}
}
