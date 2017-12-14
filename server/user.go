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
	fmt.Println("disconnect user", u.id)
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

		fmt.Printf("recv data from user(%v), len is %v\n", u.id, n)

		// 推入通道的发送数据队列
		pkg := comm.NewPacket(u.id, comm.Cmd_Data, recv[:n])
		sessionForTunnel.SendPacket(pkg)
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

func (s *Session) SendData(data []byte) {
	s.send <- data
}

func (s *Session) SendPacket(pkg *comm.Packet) {
	data := comm.EncodePacket(pkg)
	//fmt.Println("send data to tunnel", data)
	s.SendData(data)
}

func (s *Session) sendLoop() {
	for {
		data := <-s.send
		utils.WriteFull(s.conn, data)
	}
}

var sessionForTunnel *Session
