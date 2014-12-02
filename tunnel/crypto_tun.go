package tunnel

import (
	"crypto/cipher"
	"errors"
)

var InvalidBlockData error = errors.New("invalid block data")

func EncryptBlocks(b cipher.Block, bs []byte) error {
	bsize := b.BlockSize()
	if len(bs)%bsize != 0 {
		return InvalidBlockData
	}

	for i := 0; i < len(bs)/bsize; i += 1 {
		b.Encrypt(bs[i*bsize:], bs[i*bsize:])
	}
	return nil
}

func DecryptBlocks(b cipher.Block, bs []byte) error {
	bsize := b.BlockSize()
	if len(bs)%bsize != 0 {
		return InvalidBlockData
	}

	for i := 0; i < len(bs)/bsize; i += 1 {
		b.Decrypt(bs[i*bsize:], bs[i*bsize:])
	}
	return nil
}
