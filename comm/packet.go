package comm

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

const PacketLen = 4

const (
	Cmd_Connect = 1
	Cmd_Data    = 2
	Cmd_Close   = 3
)

type Packet struct {
	UserID int
	CmdID  int
	Data   []byte
}

func NewPacket(userID int, cmdID int, data []byte) *Packet {
	return &Packet{
		UserID: userID,
		CmdID:  cmdID,
		Data:   data,
	}
}

func (p *Packet) Len() int {
	return len(p.Data)
}

func Encode(pkg *Packet) []byte {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(pkg)
	if err != nil {
		fmt.Println("encode fail", err)
	}
	return network.Bytes()
}

func Decode([]byte) *Packet {
	var rv Packet
	var network bytes.Buffer
	enc := gob.NewDecoder(&network)
	err := enc.Decode(&rv)
	if err != nil {
		fmt.Println("decode fail", err)
	}
	return &rv
}
