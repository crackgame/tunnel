package comm

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"net"
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

func encode(data interface{}) []byte {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		fmt.Println("gob encode fail", err)
	}
	return buf.Bytes()
}

func decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(to)
	if err != nil {
		fmt.Println("gob decode fail", err)
	}
	return err
}

// EncodePacket 包含大小字段
func EncodePacket(pkg *Packet) []byte {
	data := encode(pkg)

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

	pkg := &Packet{}
	err = decode(bs, pkg)
	return pkg, err
}
