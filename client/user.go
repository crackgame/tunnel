package client

import (
	"fmt"
	"net"
	"os"
	"tunnel/comm"
	"tunnel/utils"
)

type User struct {
	id      int
	session *Session
}

func NewUser(userID int, conn net.Conn) *User {
	return &User{
		id:      userID,
		session: NewSession(conn),
	}
}

func (u *User) Run() {
	go u.recvLoop()
	go u.session.sendLoop()
}

func (u *User) Disconnect() {
	fmt.Println("disconnect user", u.id)
	u.session.conn.Close()
}

func (u *User) recvLoop() {
	for {
		recv := make([]byte, comm.RecvBuffSize)
		n, err := u.session.conn.Read(recv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			break
		}

		if sessionForTunnel == nil {
			fmt.Println("tunnel is not conneted")
			continue
		}

		pkg := comm.NewPacket(u.id, comm.Cmd_Data, recv[:n])
		sessionForTunnel.SendPacket(pkg)

		fmt.Println("recv data from inner, len is", n)
	}
}

//////////////////////////////
// session

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

func (s *Session) SendData(data []byte) {
	s.send <- data
}

func (s *Session) SendPacket(pkg *comm.Packet) {
	data := comm.EncodePacket(pkg)
	s.SendData(data)
}

func (s *Session) sendLoop() {
	for {
		data := <-s.send
		utils.WriteFull(s.conn, data)
	}
}

var sessionForTunnel *Session
