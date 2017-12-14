package comm

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"tunnel/comm"
	"tunnel/utils"
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

// EncodePacket 包含大小字段
func EncodePacket(pkg *Packet) []byte {
	data := Encode(pkg)

	dataLen := len(data)
	bs := make([]byte, PacketLen)
	binary.LittleEndian.PutUint32(bs, uint32(dataLen))

	data = append(bs, data...)
	return data
}

// DecodePacket 读取包
func DecodePacket(conn net.Conn) (*Packet, error) {
	// read pakcet len
	headerLen, err := utils.ReadInt32(conn)
	if err != nil {
		return nil, err
	}

	// read packet data
	bs := make([]byte, headerLen)
	_, err = io.ReadFull(conn, bs)
	if err != nil {
		return nil, err
	}
	pkg := comm.Decode(bs)
	return pkg, nil
}
