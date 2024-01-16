package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

const defaultHmacKey = "Qw12#"

func Sha256(data string, key ...string) string {
	encryptKey := defaultHmacKey
	if len(key) > 0 {
		encryptKey = key[0]
	}
	hmac := hmac.New(sha256.New, []byte(encryptKey))
	hmac.Write([]byte(data))
	return hex.EncodeToString(hmac.Sum(nil))
}
