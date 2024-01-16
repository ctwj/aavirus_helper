package util

import (
	"testing"
)

func TestAes(t *testing.T) {
	data := "abcdef"
	cc := []byte(data)
	t.Log(len(data), len(cc))
	ss, err := AesEncrypt(data)
	t.Log(ss, err)
	if nil != err {
		return
	}
	src, err := AesDecrypt(ss)
	t.Log(src, err)
}
