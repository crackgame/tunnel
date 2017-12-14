package server

import (
	"fmt"
	"net"
	"tunnel/comm"
	"tunnel/utils"
)

type User struct {
	id      int
	session *Session
}

func NewUser(conn net.Conn) *User {
	return &User{
		id:      comm.AutoIncrementID(),
		session: NewSession(conn),
	}
}

func (u *User) Run() {
	u.onConnect()
	go u.recvLoop()
	go u.session.sendLoop()
}

func (u *User) Disconnect() {
	u.session.conn.Close()
}

func (u *User) onConnect() {
	pkg := comm.NewPacket(u.id, comm.Cmd_Connect, []byte{})
	sessionForTunnel.SendPacket(pkg)
}

func (u *User) onClose() {
	pkg := comm.NewPacket(u.id, comm.Cmd_Close, []byte{})
	sessionForTunnel.SendPacket(pkg)
}

func (u *User) recvLoop() {
	for {
		recv := make([]byte, comm.RecvBuffSize)
		n, err := u.session.conn.Read(recv)
		if err != nil {
			u.onClose()
			RemoveUser(u.id)
			u.session.conn.Close()
			break
		}

		if sessionForTunnel == nil {
			fmt.Println("tunnel is not conneted")
			continue
		}

		// 推入通道的发送数据队列
		pkg := comm.NewPacket(u.id, comm.Cmd_Data, recv[:n])
		sessionForTunnel.SendPacket(pkg)

		fmt.Println("recv data from user, len is", n)
	}
}

//////////

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

func (s *Session) SendPacket(pkg *comm.Packet) {
	data := comm.Encode(pkg)
	if sessionForTunnel == nil {
		fmt.Println("tunnel is not opened", *pkg)
		return
	}

	sessionForTunnel.send <- data
}

func (s *Session) sendLoop() {
	for {
		data := <-s.send
		utils.WriteFull(s.conn, data)
	}
}

var sessionForTunnel *Session
