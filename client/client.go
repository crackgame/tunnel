package client

import (
	"fmt"
	"net"
	"tunnel/comm"
	"tunnel/utils"
)

func Run(host string, tunnelPort int, innerPort int) {
	runForTunnelClient(host, tunnelPort, innerPort)
	runForInnerClient(innerPort)
}

func runForInnerClient(port int) {
	innerAddr := fmt.Sprintf("%v:%v", "127.0.0.1", port)
	fmt.Println("connect", innerAddr)
	conn, err := net.Dial("tcp", innerAddr)
	utils.CheckError(err)

	user.sessionForInner = NewSession(conn)
	go user.sessionForInner.SendLoop()

	// 开线程跑接收
	go func() {
		for {
			recv := make([]byte, comm.RecvBuffSize)
			n, err := conn.Read(recv)
			utils.CheckError(err)

			if user.sessionForTunnel == nil {
				fmt.Println("tunnel is not conneted")
				continue
			}

			user.sessionForTunnel.send <- recv[:n]

			fmt.Println("recv data from inner, len is", n)
		}
	}()
}

func runForTunnelClient(host string, port int, innerPort int) {
	var err error
	addr := fmt.Sprintf("%v:%v", host, port)
	conn, err := net.Dial("tcp", addr)
	utils.CheckError(err)
	defer conn.Close()

	user.sessionForTunnel = NewSession(conn)
	go user.sessionForTunnel.SendLoop()

	for {
		recv := make([]byte, comm.RecvBuffSize)
		n, err := conn.Read(recv)
		utils.CheckError(err)

		// 如果收到数据，连接内网指定端口，打通通道
		if user.sessionForInner == nil {
			runForInnerClient(innerPort)
		}

		// 写入接收到的数据到本地端口
		user.sessionForInner.send <- recv[:n]

		fmt.Println("recv data from tunnel, len is", n)
	}
}
