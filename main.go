package main

import (
	"flag"
	"fmt"
	"tunnel/client"
	"tunnel/server"
)

var Mode *string = flag.String("mode", "s", "s/c")
var UserPort *int = flag.Int("up", 9983, "user port(外网端口)")
var TunnelPort *int = flag.Int("tp", 20012, "tunnel port(tcp通道端口)")
var Host *string = flag.String("host", "127.0.0.1", "服务器ip")
var InnerPort *int = flag.Int("port", 9983, "内网端口")

// exp.
// server: tunnel -mode s -up 9983 -tp 20012
// client: tunnel -mode c -host 127.0.0.1 -tp 20012 -port 9983

func main() {
	flag.Parse()
	// if flag.NFlag() < 3 {
	// 	flag.PrintDefaults()
	// 	return
	// }

	test := 1
	_ = test

	if *Mode == "s" {
		fmt.Println("run as server")
		server.Run(*UserPort, *TunnelPort)
	} else if *Mode == "c" {
		fmt.Println("run as client")
		client.Run(*Host, *TunnelPort, *InnerPort)
	}
}
