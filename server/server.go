package server

import (
	"fmt"
	"net"
	"tunnel/comm"
	"tunnel/utils"
)

func Run(userPort int, tunnelPort int) {
	go runForUser(userPort)
	runForTunnel(tunnelPort)
}

func runForUser(port int) {
	addr := fmt.Sprintf(":%v", port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	utils.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	utils.CheckError(err)
	fmt.Println("listen for user on", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		user := NewUser(conn)
		AddUser(user) // add user to cache
		user.Run()
	}
}

func runForTunnel(port int) {
	addr := fmt.Sprintf(":%v", port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr) //获取一个tcpAddr
	utils.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr) //监听一个端口
	utils.CheckError(err)
	fmt.Println("listen for tunnel on", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		sessionForTunnel = NewSession(conn)
		go sessionForTunnel.sendLoop()

		fmt.Println("client connect tunnel success!")

		for {
			// read packect
			pkg, err := comm.DecodePacket(conn)
			if err != nil {
				conn.Close()
				break
			}

			user := FindUser(pkg.UserID)
			if user == nil {
				fmt.Println("not found user", pkg.UserID)
				break
			}

			// 推入从通道的接收到的数据给用户队列
			switch pkg.CmdID {
			case comm.Cmd_Data:
				user.session.SendData(pkg.Data)
			case comm.Cmd_Close:
				user.Disconnect()
			default:
				fmt.Println("runForTunnel, unknow packet cmdID", pkg.CmdID)
			}

			fmt.Printf("recv data from tunnel user(%v), len is %v\n", user.id, pkg.Len())
		}
	}
}
