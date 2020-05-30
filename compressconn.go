package aeadconn

import (
	"crypto/cipher"
	"io"
	"net"

	"github.com/klauspost/compress/s2"
	stream "github.com/maoxs2/go-aead-iostream"
)

// AEADCompressConn uses https://github.com/klauspost/compress/tree/master/s2#s2-compression to compress the data
type AEADCompressConn struct {
	net.Conn

	compressW io.Writer
	compressR io.Reader

	pipeR *io.PipeReader
	pipeW *io.PipeWriter

	cryptoW io.Writer
	cryptoR io.Reader
}

// raw <--> 1. Encode <--> 2. Encrypt(Write) <-ciphertext-> 1. Decrypt(Read) <--> 2. Decode <--> raw
func NewAEADCompressConn(seed []byte, chunkSize int, conn net.Conn, aead cipher.AEAD) *AEADCompressConn {
	pipeR, pipeW := io.Pipe()

	return &AEADCompressConn{
		Conn: conn,

		compressW: s2.NewWriter(pipeW),
		compressR: s2.NewReader(pipeR),

		pipeR: pipeR,
		pipeW: pipeW,

		cryptoW: stream.NewStreamWriteCloser(seed, chunkSize, conn, aead),
		cryptoR: stream.NewStreamReader(seed, chunkSize, conn, aead),
	}
}

func (cc *AEADCompressConn) Close() error {
	return cc.Conn.Close()
}

func (cc *AEADCompressConn) Write(b []byte) (i int, err error) {
	i, err = cc.compressW.Write(b)
	io.Copy(cc.cryptoW, cc.pipeR)
	return
}

func (cc *AEADCompressConn) Read(b []byte) (int, error) {
	io.Copy(cc.pipeW, cc.cryptoR)
	return cc.compressR.Read(b)
}
