package common

import (
	"errors"
	"fmt"

	"golang.org/x/net/websocket"
)

type Player struct {
	User
	Conn *websocket.Conn
}

func (p Player) String() string {
	return fmt.Sprintf("Player{User: %v}", p.User)
}

func (p Player) Equals(other Player) bool {
	return p.User.Equals(other.User)
}

type PlayerQueue struct {
	list []Player
}

func NewPlayerQueue() *PlayerQueue {
	return &PlayerQueue{list: make([]Player, 0)}
}

func (q *PlayerQueue) Enqueue(item Player) {
	q.list = append(q.list, item)
}

func (q *PlayerQueue) Dequeue() (Player, error) {
	if len(q.list) == 0 {
		return Player{}, errors.New("queue is empty")
	}
	item := q.list[0]
	q.list = q.list[1:]
	return item, nil
}

func (q *PlayerQueue) Remove(player Player) (Player, error) {
	var nl []Player
	for i, v := range q.list {
		if v.Equals(player) {
			nl = append(q.list[:i], q.list[i+1:]...)
			q.list = nl
			return player, nil
		}
	}

	return Player{}, errors.New("player not in queue")
}

func (q *PlayerQueue) InQueue(player Player) bool {
	for _, v := range q.list {
		if v.Equals(player) {
			return true
		}
	}

	return false
}

func (q *PlayerQueue) IsEmpty() bool {
	return len(q.list) == 0
}

func (q *PlayerQueue) Size() int {
	return len(q.list)
}
