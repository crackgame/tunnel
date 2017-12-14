package utils

import (
	"encoding/binary"
	"io"
	"net"
)

func WriteFull(w io.Writer, buf []byte) (n int, err error) {
	var nn int
	var sends int
	bufsize := len(buf)
	for sends < bufsize && err == nil {
		nn, err = w.Write(buf[sends:])
		sends += nn
	}

	if sends > bufsize {
		err = io.ErrUnexpectedEOF
	}

	return
}

func ReadInt32(conn net.Conn) (int, error) {
	bs := make([]byte, 4)
	n, err := io.ReadFull(conn, bs)
	if err != nil {
		return n, err
	}
	headerLen := int(binary.LittleEndian.Uint16(bs))
	return headerLen, nil
}
