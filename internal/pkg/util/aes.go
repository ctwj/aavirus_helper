package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

const (
	// The key argument should be the AES key,
	// either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
	defaultAesKey = "d3bf5197115e87753350349fd3besA#!"
	defalutSalt   = "Asc!3rAsc!3r"
)

func AesEncrypt(data string, key ...string) (string, error) {
	encryptKey := defaultAesKey
	if len(key) > 0 {
		encryptKey = key[0]
	}
	block, err := aes.NewCipher([]byte(encryptKey))
	if nil != err {
		return ``, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if nil != err {
		return ``, err
	}
	ciphertext := aesgcm.Seal(nil, []byte(defalutSalt), []byte(data), nil)
	return hex.EncodeToString(ciphertext), nil
}

func AesDecrypt(data string, key ...string) (string, error) {
	encryptKey := defaultAesKey
	if len(key) > 0 {
		encryptKey = key[0]
	}
	ciphertext, err := hex.DecodeString(data)
	if nil != err {
		return ``, err
	}
	block, err := aes.NewCipher([]byte(encryptKey))
	if nil != err {
		return ``, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if nil != err {
		return ``, err
	}
	plaintext, err := aesgcm.Open(nil, []byte(defalutSalt), ciphertext, nil)
	if err != nil {
		return ``, err
	}
	return string(plaintext), nil
}
