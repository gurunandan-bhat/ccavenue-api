package native

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func Encrypt(plaintext string, bKey []byte, bIV []byte) ([]byte, error) {

	blockSize := aes.BlockSize

	bPlaintext := pkcs5Padding([]byte(plaintext), blockSize)

	block, err := aes.NewCipher(bKey)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)

	return ciphertext, nil
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {

	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(padtext, ciphertext...)
}
