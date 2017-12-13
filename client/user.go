package client

import (
	"net"
	"tunnel/utils"
)

type User struct {
	sessionForInner  *Session
	sessionForTunnel *Session
}

type Session struct {
	conn net.Conn
	recv chan []byte
	send chan []byte
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		conn: conn,
		recv: make(chan []byte),
		send: make(chan []byte),
	}
}

func (s *Session) SendLoop() {
	for {
		data := <-s.send
		utils.WriteFull(s.conn, data)
	}
}

var user User
