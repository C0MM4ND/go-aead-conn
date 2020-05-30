package aeadconn

import (
	"crypto/cipher"
	"net"

	"github.com/klauspost/compress/s2"
	stream "github.com/maoxs2/go-aead-iostream"
)

// AEADS2Conn uses https://github.com/klauspost/compress/tree/master/s2#s2-compression to compress the data
type AEADS2Conn struct {
	net.Conn
	*s2.Writer
	*s2.Reader
}

func NewAEADS2Conn(seed []byte, chunkSize int, conn net.Conn, aead cipher.AEAD) *AEADS2Conn {
	return &AEADS2Conn{
		Conn:   conn,
		Writer: s2.NewWriter(stream.NewStreamWriteCloser(seed, chunkSize, conn, aead)),
		Reader: s2.NewReader(stream.NewStreamReader(seed, chunkSize, conn, aead)),
	}
}

func (cc *AEADS2Conn) Close() error {
	return cc.Conn.Close()
}

func (cc *AEADS2Conn) Write(b []byte) (int, error) {
	return cc.Writer.Write(b)
}

func (cc *AEADS2Conn) Read(b []byte) (int, error) {
	return cc.Reader.Read(b)
}
