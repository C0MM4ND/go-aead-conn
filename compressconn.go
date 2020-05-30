package aeadconn

import (
	"crypto/cipher"
	"io"
	"net"

	snappy "github.com/klauspost/compress/snappy"
	stream "github.com/maoxs2/go-aead-iostream"
)

// AEADCompressConn uses https://github.com/klauspost/compress/tree/master/s2#s2-compression to compress the data
type AEADCompressConn struct {
	net.Conn

	w io.Writer
	r io.Reader
}

// raw <--> 1. Encode <--> 2. Encrypt(Write) <-ciphertext-> 1. Decrypt(Read) <--> 2. Decode <--> raw
func NewAEADCompressConn(seed []byte, chunkSize int, conn net.Conn, aead cipher.AEAD) *AEADCompressConn {
	cc := &AEADCompressConn{
		Conn: conn,
		w:    snappy.NewWriter(stream.NewStreamWriteCloser(seed, chunkSize, conn, aead)),
		r:    snappy.NewReader(stream.NewStreamReader(seed, chunkSize, conn, aead)),
	}

	return cc
}

func (cc *AEADCompressConn) Close() error {
	return cc.Conn.Close()
}

func (cc *AEADCompressConn) Write(b []byte) (i int, err error) {
	return cc.w.Write(b)
}

func (cc *AEADCompressConn) Read(b []byte) (int, error) {
	return cc.r.Read(b)
}
