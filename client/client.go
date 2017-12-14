package client

import (
	"fmt"
	"io"
	"net"
	"tunnel/comm"
	"tunnel/utils"
)

func Run(host string, tunnelPort int, innerPort int) {
	runForTunnelClient(host, tunnelPort, innerPort)
}

func runForInnerClient(userID int, port int) {
	innerAddr := fmt.Sprintf("%v:%v", "127.0.0.1", port)
	fmt.Println("connect", innerAddr)
	conn, err := net.Dial("tcp", innerAddr)
	utils.CheckError(err)

	user := NewUser(userID, conn)
	AddUser(user) // add user to cache
	user.Run()
}

func runForTunnelClient(host string, port int, innerPort int) {
	var err error
	addr := fmt.Sprintf("%v:%v", host, port)
	conn, err := net.Dial("tcp", addr)
	utils.CheckError(err)
	defer conn.Close()

	sessionForTunnel = NewSession(conn)
	go sessionForTunnel.sendLoop()

	for {
		// read pakcet len
		headerLen, err := utils.ReadInt32(conn)
		if err != nil {
			conn.Close()
			break
		}

		// read packet data
		bs := make([]byte, headerLen)
		_, err = io.ReadFull(conn, bs)
		if err != nil {
			conn.Close()
			break
		}
		pkg := comm.Decode(bs)

		switch pkg.CmdID {
		case comm.Cmd_Connect:
			runForInnerClient(pkg.UserID, innerPort)
		case comm.Cmd_Data:
			user := FindUser(pkg.UserID)
			if user != nil {
				user.session.SendData(pkg.Data)
			} else {
				fmt.Println("Cmd_Data not found user", pkg.UserID)
			}
		case comm.Cmd_Close:
			user := FindUser(pkg.UserID)
			if user != nil {
				user.Disconnect()
			} else {
				fmt.Println("Cmd_Close not found user", pkg.UserID)
			}
		}

		fmt.Println("recv data from tunnel, len is", pkg.Len())
	}
}
