package tunnel

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rc4"
)

type cipherMaker interface {
	// return encrypter/ decrypter
	NewCipher(key) (cipher.Block, error)
}

type CipherConfig struct {
	Name      string
	KeySize   int
	BlockSize int
	maker     cipherMaker
}

func (ctx *CipherConfig) NewCipher(key) (cipher.Block, error) {
	return ctx.maker.NewCipher(key)
}

type DESCipherMaker struct {
	is3des bool
}

func (m *DESCipherMaker) NewCipher(key) (cipher.Block, error) {
	var block cipher.Block
	var err error

	if m.is3des {
		if block, err = des.NewTripleDESCipher(key); err != nil {
			return nil, err
		}
	} else {
		if block, err = des.NewCipher(key); err != nil {
			return nil, err
		}
	}

	return block, nil
}

type AESCipherMaker struct{}

func (m *AESCipherMaker) NewCipher(key) (cipher.Block, error) {
	if block, err := aes.NewCipher(key); err != nil {
		return nil, err
	} else {
		return block, nil
	}
}

var ciphers map[string]*CipherConfig

func init() {
	ciphers = make(map[string]*CipherConfig)
	ciphers["des"] = &CipherConfig{
		Name:      "des",
		KeySize:   8,
		BlockSize: des.BlockSize,
		maker:     &DESCipherMaker{is3des: false}}
	ciphers["3des-192"] = &CipherConfig{
		Name:      "3des-192",
		KeySize:   24,
		BlockSize: des.BlockSize,
		maker:     &DESCipherMaker{is3des: true}}
	ciphers["aes-128"] = &CipherConfig{
		Name:      "aes-128",
		KeySize:   16,
		BlockSize: aes.BlockSize,
		maker:     new(AESCipherMaker)}
	ciphers["aes-192"] = &CipherConfig{
		Name:      "aes-192",
		KeySize:   24,
		BlockSize: aes.BlockSize,
		maker:     new(AESCipherMaker)}
	ciphers["aes-256"] = &CipherConfig{
		Name:      "aes-256",
		KeySize:   32,
		BlockSize: aes.BlockSize,
		maker:     new(AESCipherMaker)}
}

func GetCipherConfig(name string) *CipherConfig {
	return ciphers[name]
}
