package common

import "golang.org/x/net/websocket"

type GameList struct {
	mapByRoomId map[string]*GameRoom
	MapByConn   map[*websocket.Conn]*GameRoom
}

func NewGameList() *GameList {
	return &GameList{
		mapByRoomId: make(map[string]*GameRoom),
		MapByConn:   make(map[*websocket.Conn]*GameRoom),
	}
}

func (l *GameList) Add(gameRoom *GameRoom) {
	l.mapByRoomId[gameRoom.RoomID] = gameRoom

	conns := make([]*websocket.Conn, 0, len(gameRoom.Players))
	for c := range gameRoom.Players {
		conns = append(conns, c)
	}

	for _, c := range conns {
		l.MapByConn[c] = gameRoom
	}
}

func (l *GameList) Remove(gameRoom *GameRoom) {
	delete(l.mapByRoomId, gameRoom.RoomID)
	conns := make([]*websocket.Conn, 0, len(gameRoom.Players))
	for c := range gameRoom.Players {
		conns = append(conns, c)
	}

	for _, c := range conns {
		delete(l.MapByConn, c)
	}
}
