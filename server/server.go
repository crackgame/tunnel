package server

import (
	"fmt"
	"net"
	"tunnel/utils"
)

func Run(userPort int, tunnelPort int) {
	go runForUser(userPort)
	runForTunnel(tunnelPort)
}

func runForUser(port int) {
	addr := fmt.Sprintf(":%v", port)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr) //获取一个tcpAddr
	utils.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr) //监听一个端口
	utils.CheckError(err)
	fmt.Println("listen for user on", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		user.sessionForUser = NewSession(conn)
		go user.sessionForUser.SendLoop()

		for {
			recv := make([]byte, 10240)
			n, err := conn.Read(recv)
			if err != nil {
				conn.Close()
				continue
			}

			if user.sessionForTunnel == nil {
				fmt.Println("tunnel is not conneted")
				continue
			}

			// 推入通道的发送数据队列
			user.sessionForTunnel.send <- recv[:n]

			fmt.Println("recv data from user, len is", n)
		}
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

		user.sessionForTunnel = NewSession(conn)
		go user.sessionForTunnel.SendLoop()

		for {
			recv := make([]byte, 10240)
			n, err := conn.Read(recv)
			if err != nil {
				conn.Close()
				continue
			}

			if user.sessionForUser == nil {
				fmt.Println("tunnel is not conneted")
				continue
			}

			// 推入从通道的接收到的数据给用户队列
			user.sessionForUser.send <- recv[:n]

			fmt.Println("recv data from tunnel, len is", n)
		}
	}
}
