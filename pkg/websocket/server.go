package websocket

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns   map[*websocket.Conn]bool
	bufsize uint
}

func NewServer(bufsize uint) *Server {
	return &Server{
		conns:   make(map[*websocket.Conn]bool),
		bufsize: bufsize,
	}
}

func (s *Server) HandleWS(ws *websocket.Conn) {
	fmt.Println("New client connection: ", ws.RemoteAddr())
	s.conns[ws] = true
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, s.bufsize)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed", ws)
				s.conns[ws] = false
				break
			}
			fmt.Println("Read error: ", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))
	}
}
