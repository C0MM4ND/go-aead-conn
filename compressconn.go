package aeadconn

import (
	"compress/gzip"
	"crypto/cipher"
	"io"
	"net"

	stream "github.com/maoxs2/go-aead-iostream"
)

// AEADCompressConn uses https://github.com/klauspost/compress/tree/master/s2#s2-compression to compress the data
type AEADCompressConn struct {
	net.Conn

	compressW io.Writer
	compressR io.Reader

	writeR *io.PipeReader
	readW  *io.PipeWriter

	cryptoW io.Writer
	cryptoR io.Reader
}

// raw <--> 1. Encode <--> 2. Encrypt(Write) <-ciphertext-> 1. Decrypt(Read) <--> 2. Decode <--> raw
func NewAEADCompressConn(seed []byte, chunkSize int, conn net.Conn, aead cipher.AEAD) *AEADCompressConn {
	writeR, writeW := io.Pipe()
	readR, readW := io.Pipe()

	compressR, err := gzip.NewReader(readR)
	if err != nil {
		panic(err)
	}

	cc := &AEADCompressConn{
		Conn: conn,

		compressW: gzip.NewWriter(writeW),
		compressR: compressR,

		writeR: writeR,
		readW:  readW,

		cryptoW: stream.NewStreamWriteCloser(seed, chunkSize, conn, aead),
		cryptoR: stream.NewStreamReader(seed, chunkSize, conn, aead),
	}

	go io.Copy(cc.cryptoW, cc.writeR)
	go io.Copy(cc.readW, cc.cryptoR)

	return cc
}

func (cc *AEADCompressConn) Close() error {
	return cc.Conn.Close()
}

func (cc *AEADCompressConn) Write(b []byte) (i int, err error) {
	i, err = cc.compressW.Write(b)
	return
}

func (cc *AEADCompressConn) Read(b []byte) (int, error) {
	return cc.compressR.Read(b)
}
