package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func Encrypt(plainText, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherTextByte := gcm.Seal(
		nonce,
		nonce,
		[]byte(plainText),
		nil,
	)

	cipherText := base64.StdEncoding.EncodeToString(cipherTextByte)

	return cipherText, nil
}

func Decrypt(cipherText, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()

	cipherTextByte, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	nonce, cipherTextByte := cipherTextByte[:nonceSize], cipherTextByte[nonceSize:]
	plainTextByte, err := gcm.Open(
		nil,
		nonce,
		cipherTextByte,
		nil,
	)
	if err != nil {
		return "", err
	}

	return string(plainTextByte), nil
}
