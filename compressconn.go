package aeadconn

import (
	"crypto/cipher"
	"net"

	"github.com/klauspost/compress/s2"
	stream "github.com/maoxs2/go-aead-iostream"
)

// AEADCompressConn uses https://github.com/klauspost/compress/tree/master/s2#s2-compression to compress the data
type AEADCompressConn struct {
	net.Conn
	*stream.StreamWriteCloser
	*stream.StreamReader
}

func NewAEADCompressConn(seed []byte, chunkSize int, conn net.Conn, aead cipher.AEAD) *AEADCompressConn {
	return &AEADCompressConn{
		Conn:              conn,
		StreamWriteCloser: stream.NewStreamWriteCloser(seed, chunkSize, conn, aead),
		StreamReader:      stream.NewStreamReader(seed, chunkSize, conn, aead),
	}
}

func (cc *AEADCompressConn) Close() error {
	return cc.StreamWriteCloser.Close()
}

func (cc *AEADCompressConn) Write(b []byte) (int, error) {
	b = s2.EncodeBetter(nil, b)
	return cc.StreamWriteCloser.Write(b)
}

func (cc *AEADCompressConn) Read(b []byte) (int, error) {
	b, _ = s2.Decode(nil, b)
	return cc.StreamReader.Read(b)
}
