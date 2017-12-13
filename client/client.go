package client

import (
	"fmt"
	"net"
	"tunnel/utils"
)

func Run(host string, tunnelPort int, innerPort int) {
	runForTunnelClient(host, tunnelPort, innerPort)
	runForInnerClient(innerPort)
}

func runForInnerClient(port int) {
	// addr := fmt.Sprintf(":%v", port)
	// tcpAddr, err := net.ResolveTCPAddr("tcp4", addr) //获取一个tcpAddr
	// utils.CheckError(err)
	// listener, err := net.ListenTCP("tcp", tcpAddr) //监听一个端口
	// utils.CheckError(err)
	// fmt.Println("listen for user on", addr)
	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		continue
	// 	}

	// 	user.sessionForUser = NewSession(conn)

	// 	for {
	// 		var recv []byte = make([]byte, 10240)
	// 		n, err := conn.Read(recv)
	// 		if err != nil {
	// 			conn.Close()
	// 			continue
	// 		}

	// 		if user.sessionForTunnel == nil {
	// 			fmt.Println("tunnel is not conneted")
	// 			continue
	// 		}

	// 		fmt.Println("aaaa", n)
	// 	}
	// 	//conn.Write([]byte(daytime))
	// 	//conn.Close()
	// }
}

func runForTunnelClient(host string, port int, innerPort int) {
	var err error
	addr := fmt.Sprintf("%v:%v", host, port)
	conn, err := net.Dial("tcp", addr)
	utils.CheckError(err)
	defer conn.Close()

	var connForInner net.Conn

	for {
		var recv []byte = make([]byte, 10240)
		n, err := conn.Read(recv)
		utils.CheckError(err)

		// 如果收到数据，连接内网指定端口，打通通道
		if connForInner == nil {
			innerAddr := fmt.Sprintf("%v:%v", "127.0.0.1", innerPort)
			fmt.Println("connect", innerAddr)
			connForInner, err = net.Dial("tcp", innerAddr)
			utils.CheckError(err)
		}

		// 写入接收到的数据到本地端口
		connForInner.Write(recv[:n])

		fmt.Println("aaaa", n)
	}

	// tcpAddr, err := net.ResolveTCPAddr("tcp4", addr) //获取一个tcpAddr
	// utils.CheckError(err)
	// listener, err := net.ListenTCP("tcp", tcpAddr) //监听一个端口
	// utils.CheckError(err)
	// fmt.Println("listen for tunnel on", addr)
	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		continue
	// 	}

	// 	user.sessionForTunnel = NewSession(conn)

	// 	daytime := time.Now().String()
	// 	fmt.Println("bbbb")
	// 	conn.Write([]byte(daytime))
	// 	conn.Close()
	// }
}
