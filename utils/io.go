package utils

import "io"

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
