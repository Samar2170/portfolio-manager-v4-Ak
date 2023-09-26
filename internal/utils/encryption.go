package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

func EncryptPassword(password, passwordDecryptionKey string) (string, error) {
	block, err := aes.NewCipher([]byte(passwordDecryptionKey))
	if err != nil {
		return "", errors.New("Error during creation of cipher block:  " + err.Error())
	}
	cipherText := make([]byte, aes.BlockSize+len(password))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.New("Error during reading of random bytes:  " + err.Error())
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(password))

	return hex.EncodeToString(cipherText), nil
}
